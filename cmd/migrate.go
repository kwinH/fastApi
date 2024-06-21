package cmd

import (
	"fastApi/app/model"
	"fastApi/core/global"
	"github.com/spf13/cobra"
)

func init() {
	serverCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Initialize the database",
		Run:   migrate,
	}
	RootCmd.AddCommand(serverCmd)
}

func migrate(cmd *cobra.Command, args []string) {
	_ = global.GDB.AutoMigrate(&model.User{})
}
