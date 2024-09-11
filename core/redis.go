package core

import (
	"context"
	"fastApi/core/global"
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func RedisInit() {
	host := viper.GetString("redis.host")
	if host == "" {
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:       viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Password:   viper.GetString("redis.password"),
		DB:         viper.GetInt("redis.db"),
		MaxRetries: 1,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		panic("连接Redis不成功" + err.Error())
	}

	if viper.IsSet("telemetry") {
		client.AddHook(redisotel.NewTracingHook())
	}

	global.Redis = client
}
