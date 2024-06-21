package mq

import (
	"context"
	"fastApi/app/model"
	"fastApi/core/logger"
	"github.com/nsqio/go-nsq"
)

type SendRegisteredEmail struct {
	BaseMQ
}

func (c *SendRegisteredEmail) HandleMessage(msg *nsq.Message) error {
	return c.Handle(msg, func(ctx context.Context, data string) error {
		model.GetUser(ctx, 1)
		logger.Log(ctx).Info("ok")
		return nil
	})
}

func init() {
	MQList = append(MQList, NewSendRegisteredEmail())
}

func NewSendRegisteredEmail() *SendRegisteredEmail {
	return &SendRegisteredEmail{
		BaseMQ: BaseMQ{
			Topic: "sendRegisteredEmail",
		}}
}
