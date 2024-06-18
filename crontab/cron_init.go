package crontab

import (
	"context"
	"fastApi/core/global"
	"fastApi/core/logger"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"

	"go.uber.org/zap"
)

type InterfaceCron interface {
	getSpec() string
	getName() string
	Run(context.Context)
}

var Cron *cron.Cron
var CronList []InterfaceCron

type Logger struct {
	Log *zap.SugaredLogger
}

func (l Logger) Info(msg string, keysAndValues ...interface{}) {
	l.Log.Debug(msg, keysAndValues)
}

func (l Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.Log.Error(err, msg, keysAndValues)
}

func CronInit() {
	newlog := Logger{
		Log: global.Log.Sugar(),
	}
	Cron = cron.New(
		cron.WithChain(
			cron.Recover(newlog),
			cron.DelayIfStillRunning(newlog),
			cron.SkipIfStillRunning(newlog),
		),
		cron.WithSeconds(),
		cron.WithLogger(newlog),
	)

	Schedule()

	Cron.Start()
}

func WithRequestId(name, traceId string) {
	logger.With(
		zap.String("traceId", traceId),
		zap.String("name", name),
	)
}

func BaseCronFuc(name string, cmd func(context.Context)) func() {
	return func() {
		traceId := uuid.New().String()
		ctx := context.WithValue(context.Background(), logger.TraceId, traceId)

		WithRequestId(name, traceId)
		cmd(ctx)
	}
}

func AddCron(cmd InterfaceCron) {
	Cron.AddFunc(
		cmd.getSpec(),
		BaseCronFuc(
			cmd.getName(),
			func(ctx context.Context) {
				cmd.Run(ctx)
			}),
	)
}
