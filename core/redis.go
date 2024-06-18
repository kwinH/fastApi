package core

import (
	"fastApi/core/global"
	"github.com/go-redis/redis"
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

	_, err := client.Ping().Result()

	if err != nil {
		panic("连接Redis不成功" + err.Error())
	}

	global.Redis = client
}
