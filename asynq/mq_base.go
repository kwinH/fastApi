package asynq

import (
	"context"
	"encoding/json"
	"fastApi/core/logger"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

var AsynqClient *asynq.Client
var MQList []InterfaceMQ

type InterfaceMQ interface {
	Producer(ctx context.Context, message []byte, opts ...asynq.Option) error
	Consumer(ctx context.Context, t *asynq.Task) error
	GetTypename() string
}

type BaseMQ struct {
	Typename string
}

func (c *BaseMQ) GetTypename() string {
	return c.Typename
}

func GetAsynqClient() {
	if AsynqClient == nil {
		r := asynq.RedisClientOpt{
			Addr:     viper.GetString("asynq.addr"),
			Password: viper.GetString("asynq.password"),
			DB:       viper.GetInt("asynq.db"),
		}
		AsynqClient = asynq.NewClient(r)
	}
}

func (c *BaseMQ) Producer(ctx context.Context, payload []byte, opts ...asynq.Option) (err error) {
	GetAsynqClient()

	traceId := ctx.Value(logger.TraceId).(string)
	if traceId == "" {
		traceId, _ = logger.CalcTraceId(ctx)
	}

	data := map[string]string{
		"traceId": traceId,
		"payload": string(payload),
	}
	payload, _ = json.Marshal(data)

	task := asynq.NewTask(c.GetTypename(), payload)

	//延时消费 5秒后消费
	//opts = append(opts, asynq.ProcessIn(5*time.Second))

	// 10秒超时
	//opts = append(opts, asynq.Timeout(10*time.Second))

	//最多重试次参数
	opts = append(opts, asynq.MaxRetry(100))

	//20秒后过期
	//opts = append(opts, asynq.Deadline(time.Now().Add(20*time.Second)))

	res, err := AsynqClient.Enqueue(task, opts...)

	if err != nil {
		logger.SLog(ctx).Errorf("加入队列(%s)失败:%v", res.Type, err.Error())
	} else {
		logger.SLog(ctx).Infof("加入队列：%s成功，参数：%s", res.Type, string(payload))
	}

	return
}
