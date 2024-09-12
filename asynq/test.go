package asynq

import (
	"context"
	"fastApi/core/logger"
	"github.com/hibiken/asynq"
)

func init() {
	MQList = append(MQList, NewTest())
}

type Test struct {
	BaseMQ
}

func NewTest() *Test {
	return &Test{
		BaseMQ: BaseMQ{
			Typename: "test",
		}}
}

func (c *Test) Consumer(ctx context.Context, t *asynq.Task) error {
	logger.SLog(ctx).Infof("type: %v, payload: %s", t.Type(), string(t.Payload()))

	return nil
}
