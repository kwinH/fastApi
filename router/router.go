package router

import (
	"fastApi/core/middleware"
	_ "fastApi/docs" // 千万不要忘了导入把你上一步生成的docs
	"fastApi/mq"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"time"
)

// init router
func NewRouter() *gin.Engine {
	r := initGin()
	loadRoute(r)
	return r
}

// init Gin
func initGin() *gin.Engine {
	//设置gin模式
	gin.SetMode(viper.GetString("api.mode"))
	engine := gin.New()

	if viper.IsSet("telemetry") {
		engine.Use(otelgin.Middleware("gin"))
	}

	engine.Use(middleware.AddTraceId())
	engine.Use(middleware.GinZap([]string{}), middleware.RecoveryWithZap(true))

	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	engine.GET("/queue_test", func(c *gin.Context) {
		err := mq.NewSendRegisteredEmail().Producer(c, []byte("test"), 1*time.Second)
		fmt.Printf("\n\n%#v\n\n", err)
		c.String(200, "pong")
	})

	return engine
}

// 加载路由
func loadRoute(r *gin.Engine) {
	apiRoute(r)
	swaggerRoute(r)
}
