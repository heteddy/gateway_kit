// @Author : detaohe
// @File   : main.go
// @Description:
// @Date   : 2022/4/26 8:31 PM

package main

import (
	"fmt"
	"gateway_kit/config"
	"gateway_kit/core/gateway"
	"gateway_kit/core/lb"
	"gateway_kit/transport"
	transportHttp "gateway_kit/transport/http"
	"github.com/lukesampson/figlet/figletlib"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

var (
	Version    bool
	Log        string
	Daemon     bool
	BannerFont string
	ConfigFile string
)

func init() {
	startCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", "./", "配置文件路径")
	startCmd.PersistentFlags().StringVarP(&Log, "log", "l", "./", "日志文件路径")
	startCmd.PersistentFlags().BoolVarP(&Daemon, "daemon", "d", false, "daemon方式启动")
	startCmd.PersistentFlags().StringVarP(&BannerFont, "fonts", "f", "./", "字体")
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
		printInfo += "\n starting...\n"
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

// CurrentAbPathByCaller 获取当前执行文件绝对路径（go run）
func CurrentAbPathByCaller(skip int) string {
	var abPath string
	_, filename, _, ok := runtime.Caller(skip)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
func printBanner(svc string) {
	//cwd, _ := os.Getwd()
	//fmt.Println("path=", cwd, "path2=", CurrentAbPathByCaller(0))
	fontsDir := filepath.Join(CurrentAbPathByCaller(0), "fonts")
	f, err := figletlib.GetFontByName(fontsDir, "standard")
	if err != nil {
		fmt.Println("Could not find that font!")
		return
	}
	fmt.Println("******************************************************************")
	figletlib.FPrintMsg(os.Stdout, " * "+svc+" * ", f, 100, f.Settings(), "left")
	fmt.Println("******************************************************************")
	fmt.Println("\n\n")
}

//func printBanner2() {
//	fmt.Println("                   _                                 _    _ _   ")
//	fmt.Println("        __ _  __ _| |_ _____      ____ _ _   _      | | _(_) |_ ")
//	fmt.Println("       / _` |/ _` | __/ _ \\ \\ /\\ / / _` | | | |_____| |/ / | __|")
//	fmt.Println("      | (_| | (_| | ||  __/\\ V  V / (_| | |_| |_____|   <| | |_ ")
//	fmt.Println("       \\__, |\\__,_|\\__\\___| \\_/\\_/ \\__,_|\\__, |     |_|\\_\\_|\\__|")
//	fmt.Println("       |___/                             |___/                  ")
//}

func startServer() {
	printBanner(config.All.Name)
	errC := make(chan error)
	go transport.Interrupt(errC)
	handler := transportHttp.MakeHttpHandler()
	svr := &http.Server{
		Addr:           ":" + config.All.Server.HttpPort,
		Handler:        handler,
		MaxHeaderBytes: 1 << 25,
	}
	go func() {
		err := svr.ListenAndServe() //改成
		errC <- err
	}()
	builder := gateway.NewProxyBuilder()
	balancer := lb.NewRoundRobin()
	// 清理
	go func() {
		// 通过配置获取负载均衡策略
		// note 暂时hardcode
		reverse := builder.BuildHttpProxy(balancer)
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
