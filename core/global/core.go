package global

import (
	"context"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const DBKey = "DB"

var (
	Trans    ut.Translator // 定义一个全局翻译器T
	Log      *zap.Logger
	SLog     *zap.SugaredLogger
	GDB      *gorm.DB // DB 数据库链接单例
	Redis    *redis.Client
	Producer *nsq.Producer
)

func DB(ctx context.Context) *gorm.DB {
	if c, ok := ctx.(*gin.Context); ok {
		return c.Value(DBKey).(*gorm.DB).WithContext(c.Request.Context())
	} else {
		return ctx.Value(DBKey).(*gorm.DB).WithContext(ctx)
	}
}
