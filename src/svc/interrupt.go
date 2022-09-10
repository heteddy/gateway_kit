// @Author : detaohe
// @File   : interrupt.go
// @Description:
// @Date   : 2022/9/7 21:21

package svc

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Interrupt(errC chan<- error) {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)
	err := fmt.Errorf("Interrupt %serviceAddrs", <-signalC)
	errC <- err
}
