package server

import (
	"fastApi/core"
	"fastApi/router"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"syscall"
)

var pidFilePath = "fastApi.pid" // 指定 pid 文件路径

var StartApi = &cobra.Command{
	Use:   "server",
	Short: "start Api Server",
	Long:  "bin/fastApi server -c config.yaml",
	Run:   serverStart,
}

func serverStart(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "stop":
			if err := stopServer(); err != nil {
				log.Fatalf("stop: %#v\n", err)
			}
			return
		case "restart":
			if err := restartServer(); err != nil {
				log.Fatalf("restart: %#v\n", err)
			}
			return
		}
	}

	if err := core.InitTrans("zh"); err != nil {
		panic(fmt.Sprintf("init trans failed, err:%v\n", err))
	}

	// 装载路由
	r := router.NewRouter()

	//	r.Run(viper.GetString("api.port"))
	s := endless.NewServer(viper.GetString("api.port"), r)

	s.BeforeBegin = func(addr string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	// 记录 pid 到文件
	if err := recordPID(); err != nil {
		log.Fatalf("Failed to write PID file: %v", err)
	}

	log.Println("Starting server...")
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// 记录 pid 到指定文件
func recordPID() error {
	pid := syscall.Getpid()
	file, err := os.Create(pidFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d\n", pid)
	if err != nil {
		return err
	}
	return nil
}

// 获取 PID
func getPID() (int, error) {
	file, err := os.Open(pidFilePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var pid int
	_, err = fmt.Fscanf(file, "%d", &pid)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// 关闭进程
func stopServer() error {
	fmt.Println("Stopping application...")
	pid, err := getPID()
	if err != nil {
		return err
	}

	// 发送 SIGTERM 信号优雅地关闭进程
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return err
	}

	log.Println("Server stopped gracefully")
	removePIDFile() // 确保在退出时删除 PID 文件

	return nil
}

func restartServer() error {
	fmt.Println("Restarting application...")
	pid, err := getPID()
	if err != nil {
		return err
	}

	if err := syscall.Kill(pid, syscall.SIGHUP); err != nil {
		return err
	}
	return nil
}

func removePIDFile() {
	if err := os.Remove(pidFilePath); err != nil {
		log.Fatalf("Failed to remove PID file: %v", err)
	}
}
