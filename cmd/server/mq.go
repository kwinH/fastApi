package server

import (
	"fastApi/mq"
	_ "fastApi/mq"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

var StartMQ = &cobra.Command{
	Use:   "mq",
	Short: "启动消息队列消费者",
	Long:  "bin/fastApi mq -c config/settings.yml",
	Run:   mqStart,
}

func mqStart(cmd *cobra.Command, args []string) {
	address := viper.GetString("nsq.consumer")
	for _, mq := range mq.MQList {
		initConsumer(mq.GetTopic(), mq.GetChannel(), address, mq)
	}
	select {}
}

// 初始化消费者
func initConsumer(topic string, channel string, address string, handler nsq.Handler) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second          //设置重连时间
	c, err := nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		panic(err)
	}
	c.SetLoggerLevel(nsq.LogLevelWarning) //支显示警告级别及以上的日志
	c.AddHandler(handler)                 // 添加消费者接口

	//建立NSQLookupd连接
	if err := c.ConnectToNSQLookupd(address); err != nil {
		panic(err)
	}
}
