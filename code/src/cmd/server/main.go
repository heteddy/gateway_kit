// @Author : detaohe
// @File   : main.go
// @Description:
// @Date   : 2022/4/26 8:31 PM

package main

import (
	"fmt"
	"gateway_kit/transport"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	errc := make(chan error)

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		err := fmt.Errorf("%s", <-sigc)
		errc <- err
	}()

	go func() {
		svr := transport.MakeHttpHandler()
		errc <- http.ListenAndServe(":9001", svr)
	}()
	<-errc
}
