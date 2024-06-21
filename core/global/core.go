package global

import (
	"context"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const DBKey = "DB"

var (
	Trans    ut.Translator // 定义一个全局翻译器T
	Log      *zap.Logger
	GDB      *gorm.DB // DB 数据库链接单例
	Redis    *redis.Client
	Producer *nsq.Producer
)

func DB(ctx context.Context) *gorm.DB {
	return ctx.Value(DBKey).(*gorm.DB)
}
