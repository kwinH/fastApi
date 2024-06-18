package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "version",
	Short: "version 子命令.",
	Long:  "这是一个version 子命令",
	Run:   runVersion,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println("version is 1.0.0")
}
