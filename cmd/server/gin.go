package server

import (
	"fastApi/core"
	"fastApi/router"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var StartApi = &cobra.Command{
	Use:   "server",
	Short: "start Api Server",
	Long:  "bin/fastApi server -c config.yaml",
	Run:   serverStart,
}

func serverStart(cmd *cobra.Command, args []string) {
	if err := core.InitTrans("zh"); err != nil {
		panic(fmt.Sprintf("init trans failed, err:%v\n", err))
	}

	// 装载路由
	r := router.NewRouter()

	s := endless.NewServer(viper.GetString("api.port"), r)
	err := s.ListenAndServe()
	if err != nil {
		log.Printf("server err: %v", err)
	}
	//	r.Run(viper.GetString("api.port"))
}
