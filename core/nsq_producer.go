package core

import (
	"fastApi/core/global"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

func NsqProducerInit() {
	// 初始化生产者
	addr := viper.GetString("nsq.producer")

	if addr == "" {
		return
	}

	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		panic(err)
	}

	err = producer.Ping()
	if err != nil {
		panic(err)
	}

	global.Producer = producer
}
