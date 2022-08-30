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

func main() {
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
