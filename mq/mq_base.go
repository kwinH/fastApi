package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fastApi/core/global"
	"fastApi/core/logger"
	"fastApi/util"
	"github.com/nsqio/go-nsq"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

var MQList []InterfaceMQ

type HandleFunc func(context.Context, string) error

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
		traceId, _ = logger.CalcTraceId(ctx)
	}

	data := map[string]string{
		"traceId": traceId,
		"message": string(message),
	}

	spanId := ctx.Value(logger.SpanId).(string)
	if spanId != "" {
		data["spanId"] = spanId
	}

	message, _ = json.Marshal(data)

	if len(delay) == 0 {
		err = producer.Publish(b.GetTopic(), message) // 发布消息
	} else {
		err = producer.DeferredPublish(b.GetTopic(), delay[0], message)
	}

	return err
}

func (b *BaseMQ) Handle(msg *nsq.Message, h HandleFunc) error {
	var log *zap.Logger
	ctx := context.Background()

	defer func() {
		if r := recover(); r != nil {
			if log == nil {
				log = logger.Log(ctx).With(
					zap.String("url", b.GetTopic()),
				)
			}
			log.Sugar().Errorf("panic: %v", r)
		}
	}()

	var span oteltrace.Span
	ctx, span, _ = util.ContextWithSpanContext(ctx, "", "", "queue-consumer", b.GetTopic())
	if span != nil {
		defer span.End()
	}

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

	if span != nil {
		span.SetAttributes(attribute.String("traceId", data["traceId"]))
	}

	ctx = context.WithValue(ctx, logger.TraceId, data["traceId"])
	ctx = logger.WithC(
		ctx,
		zap.String("traceId", data["traceId"]),
	)

	log = logger.Log(ctx).With(
		zap.String("url", b.GetTopic()),
		zap.String("params", data["message"]),
		zap.Uint16("attempts", msg.Attempts),
	)

	err = h(ctx, data["message"])

	endTime := time.Now()
	latencyTime := endTime.Sub(startTime)
	log.With(
		zap.Duration("runtime", latencyTime),
	)
	if err != nil {
		log.Error("任务执行失败： " + err.Error())
	} else {
		log.Info("任务执行成功")
	}

	return err
}
