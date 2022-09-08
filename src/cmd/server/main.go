// @Author : detaohe
// @File   : main.go
// @Description:
// @Date   : 2022/4/26 8:31 PM

package main

import (
	"fmt"
	"gateway_kit/config"
	"gateway_kit/gateway"
	"gateway_kit/lb"
	"gateway_kit/svc"
	"gateway_kit/transport"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

var (
	Version    bool
	Log        string
	Daemon     bool
	ConfigFile string
)

func init() {
	startCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", "./", "配置文件路径")
	startCmd.PersistentFlags().StringVarP(&Log, "log", "l", "./", "日志文件路径")
	startCmd.PersistentFlags().BoolVarP(&Daemon, "daemon", "d", false, "daemon方式启动")
	rootCmd.PersistentFlags().BoolVarP(&Version, "version", "v", false, "版本信息")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
}

var rootCmd = &cobra.Command{
	Use:   "gateway-kit",
	Short: "api gateway tool kit",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if Version {
			ret := "\t BUILD  TIME:" + config.BuildTime + "\n"
			ret += "\t GIT VERSION:" + config.GitVersion + "\n"
			ret += "\t GO  VERSION:" + config.GoVersion + "\n"
			fmt.Println(ret)
		} else {
			fmt.Println("输入-h查看帮助")
		}
	},
}
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动服务",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if Daemon {
			cmd := exec.Command("./gateway-kit", "startServer")
			if err := cmd.Start(); err != nil {
				panic(err)
			}
			fmt.Printf("starting, [pid] %d running...\n", cmd.Process.Pid)
			if err := ioutil.WriteFile("pid.lock", []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0666); err != nil {
				panic(err)
			}
		}
		if err := config.InitConfigure(ConfigFile, Log); err != nil {
			panic(err)
		}
		var printInfo string
		printInfo += "Info:\n"
		printInfo += "\t BUILD  TIME:" + config.BuildTime + "\n"
		printInfo += "\t GIT VERSION:" + config.GitVersion + "\n"
		printInfo += "\t GO  VERSION:" + config.GoVersion + "\n"
		printInfo += "\n starting..."
		fmt.Printf(printInfo)
		startServer()
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止服务",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func startServer() {
	errC := make(chan error)
	go svc.Interrupt(errC)
	handler := transport.MakeHttpHandler()
	svr := &http.Server{
		Addr:           ":" + config.All.Http.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 25,
	}
	go func() {
		err := svr.ListenAndServe() //改成
		errC <- err
	}()
	// 清理
	go func() {
		// 通过配置获取负载均衡策略
		// note 暂时hardcode
		reverse := gateway.NewHttpReverseProxy(svc.NewServiceRepo(), lb.NewRandomLB())
		proxySvr := &http.Server{
			Addr:           ":" + config.All.Gateway.HttpPort,
			Handler:        reverse,
			MaxHeaderBytes: 1 << 25,
		}
		log.Fatalln(proxySvr.ListenAndServe())
	}()
	<-errC
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
