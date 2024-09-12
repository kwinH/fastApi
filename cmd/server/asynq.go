package server

import (
	"context"
	"encoding/json"
	nq "fastApi/asynq"
	"fastApi/core/global"
	"fastApi/core/logger"
	"fastApi/util"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	StartAsynq = &cobra.Command{
		Use:     "asynq",
		Short:   "Welcome to a Tour of Asynq!",
		Example: "Asynq 是一个由Redis实现的异步队列 Go 库",
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() error {
	// asynq server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     viper.GetString("asynq.addr"),
			Password: viper.GetString("asynq.password"),
			DB:       viper.GetInt("asynq.db"),
		},
		asynq.Config{
			Logger: global.SLog,
			//			Concurrency: 20, //worker数量,默认启动的worker数量是服务器的CPU个数
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				delay := 1 * time.Second
				if n < 15 {
					delay = 60 * time.Second
				} else if n < 25 {
					delay = 300 * time.Second
				} else if n < 60 {
					delay = 600 * time.Second
				} else {
					delay = 3600 * time.Second
				}

				return delay
			},
		},
	)

	mux := asynq.NewServeMux()

	// some middlewares
	mux.Use(func(next asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			var log *zap.Logger
			var span oteltrace.Span

			defer func() {
				if r := recover(); r != nil {
					if log == nil {
						log = logger.Log(ctx).With(
							zap.String("url", t.Type()),
						)
					}
					log.Sugar().Errorf("panic: %v", r)
				}
			}()

			ctx, span, _ = util.ContextWithSpanContext(ctx, "", "", "queue-consumer", t.Type())
			if span != nil {
				defer span.End()
			}

			startTime := time.Now()
			var data map[string]string
			err := json.Unmarshal(t.Payload(), &data)

			if err != nil {
				log = logger.Log(ctx).With(
					zap.String("url", t.Type()),
					zap.String("params", string(t.Payload())),
				)

				log.Error("数据解析失败： " + err.Error())
			}

			if span != nil {
				span.SetAttributes(attribute.String("traceId", data["traceId"]))
			}

			ctx = logger.WithC(ctx,
				zap.String("url", t.Type()),
				zap.String("traceId", data["traceId"]),
			)

			t = asynq.NewTask(t.Type(), []byte(data["payload"]))
			err = next.ProcessTask(ctx, t)

			endTime := time.Now()
			latencyTime := endTime.Sub(startTime)
			log = logger.Log(ctx).With(
				zap.String("params", data["payload"]),
				zap.Duration("runtime", latencyTime),
			)
			if err != nil {
				log.Error("任务执行失败： " + err.Error())
			} else {
				log.Info("任务执行成功")
			}

			return err
		})
	})

	for _, q := range nq.MQList {
		mux.HandleFunc(q.GetTypename(), q.Consumer)
	}

	// start server
	if err := srv.Start(mux); err != nil {
		panic(err)
	}

	// Wait for termination signal.
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT,
		//syscall.SIGUSR1, syscall.SIGUSR2,
	)

	for {
		s := <-c
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("Program Exit...", s)
			srv.Shutdown()
			srv.Stop()
			return nil
		//case syscall.SIGUSR1:
		//	fmt.Println("usr1 signal", s)
		//case syscall.SIGUSR2:
		//	fmt.Println("usr2 signal", s)
		default:
			fmt.Println("other signal", s)
		}
	}
}
