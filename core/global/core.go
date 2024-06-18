package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Trans    ut.Translator // 定义一个全局翻译器T
	Log      *zap.Logger
	SLog     *zap.SugaredLogger
	DB       *gorm.DB // DB 数据库链接单例
	Redis    *redis.Client
	Producer *nsq.Producer
)
