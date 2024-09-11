package core

import (
	"fastApi/core/global"
	logger2 "fastApi/core/logger"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

// Database 在中间件中初始化mysql链接
func Database() {
	host := viper.GetString("database.host")
	if host == "" {
		return
	}

	newLogger := logger2.NewGormLog(global.Log)
	db, err := gorm.Open(mysql.Open(
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			viper.GetString("database.username"),
			viper.GetString("database.password"),
			host,
			viper.GetString("database.port"),
			viper.GetString("database.db"),
			viper.GetString("database.charset"),
		)), &gorm.Config{
		Logger: newLogger,
	})
	// Error
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if viper.IsSet("telemetry") {
		if err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
			panic(err)
		}
	}

	//设置连接池
	//空闲
	sqlDB.SetMaxIdleConns(viper.GetInt("database.max_idle_conn"))
	//打开
	sqlDB.SetMaxOpenConns(viper.GetInt("database.max_open_conn"))
	global.GDB = db
}
