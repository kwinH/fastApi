package crontab

import (
	"context"
	"fastApi/app/model"
	"fastApi/core/logger"
)

func init() {
	CronList = append(CronList, testJob{})
}

type testJob struct {
}

func (j testJob) getSpec() string {
	return "@every 3s"
}

func (j testJob) getName() string {
	return "test1"
}

func (j testJob) Run(ctx context.Context) {
	model.GetUser(ctx, 1)
	logger.Log(ctx).Info("tick every 1 second run once")
}
