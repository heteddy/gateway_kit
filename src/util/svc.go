// @Author : detaohe
// @File   : svc_name.go
// @Description:
// @Date   : 2022/9/11 14:38

package util

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Svc struct {
	wg        sync.WaitGroup
	ctx       context.Context
	cancelCtx context.CancelFunc
	name      string
}

type SvcEndpoint func()

func NewSvc(name string) *Svc {
	_ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(_ctx)
	return &Svc{
		wg:        sync.WaitGroup{},
		ctx:       ctx,
		cancelCtx: cancelCtx,
		name:      name,
	}
}

func (s *Svc) wait() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigChan:
		s.Stop()
	case <-s.ctx.Done():
		//
	}
	s.wg.Wait()
}

func (s *Svc) Start(endpoint SvcEndpoint) {
	s.wg.Add(1)

	go func() {
		defer func() {
			s.wg.Done()
		}()
		endpoint()
	}()

	s.wait()
}

func (s *Svc) Stop() {
	s.cancelCtx()
}

type TickerSvc struct {
	Svc
	ticker      *time.Ticker
	triggerChan chan struct{}
	runAtOnce   bool
}

func NewTickerSvc(name string, d time.Duration, atOnce bool) *TickerSvc {
	return &TickerSvc{
		Svc:         *NewSvc(name),
		ticker:      time.NewTicker(d),
		runAtOnce:   atOnce,
		triggerChan: make(chan struct{}),
	}
}

func (s *TickerSvc) Trigger() {
	s.triggerChan <- struct{}{}
}

func (s *TickerSvc) Start(endpoint SvcEndpoint) {
	fmt.Printf("task name: %s starting...\n", s.Svc.name)
	//if s.runAtOnce {
	//	s.wg.Add(1)
	//	go func() {
	//		defer func() {
	//			s.wg.Done()
	//		}()
	//		endpoint()
	//	}()
	//
	//}
	s.wg.Add(1)
	go func() {
		defer func() {
			s.wg.Done()
			s.ticker.Stop()
			fmt.Printf("task name: %s exit...\n", s.Svc.name)
		}()
		fmt.Printf("task name: %s running...\n", s.Svc.name)
	loop:
		for {
			select {
			case <-s.ctx.Done():
				break loop
			case <-s.triggerChan:
				endpoint()
			case <-s.ticker.C:
				endpoint()
			}
		}
	}()
	if s.runAtOnce {
		fmt.Printf("立即执行 task name: %s\n", s.Svc.name)
		s.Trigger()
	}
	s.wait()
}
