// @Author : detaohe
// @File   : main.go
// @Description:
// @Date   : 2022/4/26 8:31 PM

package main

import (
	"fmt"
	"gateway_kit/gateway"
	"gateway_kit/lb"
	"gateway_kit/svc"
	"gateway_kit/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//func printBanner() {
//	cwd, _ := os.Getwd()
//	fontsdir := filepath.Join(cwd, "fonts")
//	f, err := figletlib.GetFontByName(fontsdir, "standard")
//	if err != nil {
//		fmt.Fprintln(os.Stderr, "Could not find that font!")
//		return
//	}
//
//	figletlib.FPrintMsg(os.Stdout, "gateway-kit", f, 100, f.Settings(), "center")
//}
func printBanner() {
	fmt.Println("                             _                                 _    _ _   ")
	fmt.Println("                  __ _  __ _| |_ _____      ____ _ _   _      | | _(_) |_ ")
	fmt.Println("                 / _` |/ _` | __/ _ \\ \\ /\\ / / _` | | | |_____| |/ / | __|")
	fmt.Println("                | (_| | (_| | ||  __/\\ V  V / (_| | |_| |_____|   <| | |_ ")
	fmt.Println("                 \\__, |\\__,_|\\__\\___| \\_/\\_/ \\__,_|\\__, |     |_|\\_\\_|\\__|")
	fmt.Println("                 |___/                             |___/                  ")
}

func main() {
	printBanner()
	errc := make(chan error)

	go func() {
		signalC := make(chan os.Signal, 1)
		signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)
		err := fmt.Errorf("%s", <-signalC)
		errc <- err
	}()

	go func() {
		svr := transport.MakeHttpHandler()
		errc <- http.ListenAndServe(":9901", svr)
	}()

	go func() {
		// 通过配置获取负载均衡策略
		// note 暂时hardcode
		reverse := gateway.NewHttpReverseProxy(svc.NewServiceRepo(), lb.NewRandomLB())
		log.Fatalln(http.ListenAndServe(":9902", reverse))
	}()

	<-errc
}
