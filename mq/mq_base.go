package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fastApi/core/global"
	"fastApi/core/logger"
	"fmt"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"time"
)

var MQList []InterfaceMQ

type InterfaceMQ interface {
	Producer(ctx context.Context, message []byte, delay ...time.Duration) error
	HandleMessage(msg *nsq.Message) error
	GetTopic() string
	GetChannel() string
}

type BaseMQ struct {
	Topic   string
	Channel string
}

func (b *BaseMQ) GetChannel() string {
	if b.Channel == "" {
		return "channel1"
	}
	return b.Channel
}

func (b *BaseMQ) GetTopic() string {
	return b.Topic
}

func (b *BaseMQ) Producer(ctx context.Context, message []byte, delay ...time.Duration) (err error) {
	producer := global.Producer

	if producer == nil {
		return errors.New("producer is nil")
	}

	traceId := ctx.Value(logger.TraceId).(string)
	if traceId == "" {
		traceId = logger.CalcTraceId()
	}

	data := map[string]string{
		"traceId": traceId,
		"message": string(message),
	}
	message, _ = json.Marshal(data)

	fmt.Printf("traceId: %s, message: %s\n", b.GetTopic(), string(message))
	if len(delay) == 0 {
		err = producer.Publish(b.GetTopic(), message) // 发布消息
	} else {
		err = producer.DeferredPublish(b.GetTopic(), delay[0], message)
	}

	return err
}

func (b *BaseMQ) Handle(msg *nsq.Message, h func(string) error) error {
	startTime := time.Now()

	var data map[string]string
	err := json.Unmarshal(msg.Body, &data)

	if err != nil {
		global.Log.With(
			zap.String("url", b.GetTopic()),
			zap.String("params", string(msg.Body)),
			zap.Uint16("attempts", msg.Attempts),
		).Error("数据解析失败： " + err.Error())
	}

	logger.With(
		zap.String("traceId", data["traceId"]),
	)
	err = h(data["message"])

	endTime := time.Now()
	latencyTime := endTime.Sub(startTime)
	log := global.Log.With(
		zap.String("url", b.GetTopic()),
		zap.String("params", data["message"]),
		zap.Uint16("attempts", msg.Attempts),
		zap.Duration("runtime", latencyTime),
	)
	if err != nil {
		log.Error("任务执行失败： " + err.Error())
	} else {
		log.Info("任务执行成功")
	}

	return err
}
