package cmd

import (
	"fastApi/cmd/server"
	"fastApi/core"
	"fmt"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:              "",
	Short:            "这是 cobra 测试程序主入口",
	Long:             `这是 cobra 测试程序主入口， 无参数启动时进入`,
	PersistentPreRun: persistentPreRun,
	Run:              runRoot,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&core.ConfigFilePath, "config", "c", "", "Start server with provided configuration file")
	RootCmd.AddCommand(server.StartApi)
	RootCmd.AddCommand(server.StartCmd)
	RootCmd.AddCommand(server.StartMQ)
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	core.CortInit()
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	fmt.Print("欢迎使用gin脚手架fastApi 可以使用 -h 查看命令")
}
