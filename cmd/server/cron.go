package server

import (
	"fastApi/crontab"
	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "cron",
	Short: "启动定时任务",
	Long:  "bin/fastApi cron -c config/settings.yml",
	Run:   crontabStart,
}

func crontabStart(cmd *cobra.Command, args []string) {
	crontab.CronInit()

	select {}
}
