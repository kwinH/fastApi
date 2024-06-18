package core

import (
	"fmt"
	"github.com/spf13/viper"
)
var ConfigFilePath string
func ViperInit() {
	if ConfigFilePath == "" { // 优先级: 命令行 > 环境变量 > 默认值
		ConfigFilePath = "./config.yaml"
		fmt.Printf("您正在使用ConfigFilePath的默认值,ConfigFilePath的路径为%v\n", ConfigFilePath)
	} else {
		fmt.Printf("您正在使用命令行的-c参数传递的值,ConfigFilePath的路径为%v\n", ConfigFilePath)
	}

	viper.SetConfigFile(ConfigFilePath)
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error ConfigFilePath file: %s \n", err))
	}
}
