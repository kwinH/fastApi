package mq

import (
	"fastApi/app/model"
	"fastApi/core/global"
	"github.com/nsqio/go-nsq"
)

type SendRegisteredEmail struct {
	BaseMQ
}

func (c *SendRegisteredEmail) HandleMessage(msg *nsq.Message) error {
	return c.Handle(msg, func(data string) error {
		model.GetUser(1)
		global.Log.Info("ok")
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
