// @Author : detaohe
// @File   : gateway_info.go
// @Description:
// @Date   : 2022/9/13 15:42

package gateway

import (
	"gateway_kit/dao"
	"sync"
)

var onceGwCtrl sync.Once
var GwController *GWController

type GWController struct {
	mutex      sync.RWMutex
	stopC      chan struct{}
	gwChan     chan *dao.GwEvent
	accessChan chan<- *AccessConfigEvt
}

func NewGwController(accessChan chan<- *AccessConfigEvt) *GWController {
	onceGwCtrl.Do(func() {
		GwController = &GWController{
			mutex:      sync.RWMutex{},
			stopC:      make(chan struct{}),
			gwChan:     make(chan *dao.GwEvent),
			accessChan: accessChan,
		}
	})
	return GwController
}

func (ctrl *GWController) In() chan<- *dao.GwEvent {
	return ctrl.gwChan
}
func (ctrl *GWController) Start() {
	go func() {
	loop:
		for {
			select {
			case <-ctrl.stopC:
				break loop
			case event, ok := <-ctrl.gwChan:
				if !ok {
					break loop
				}
				entity := event.Entity
				accessConfig := &AccessConfigEvt{
					EventType: event.EventType,
					BlockIP:   entity.BlockList,
					AllowIP:   entity.AllowList,
					Name:      entity.Name,
					Category:  ACCESS_CONTROL_GATEWAY,
				}
				ctrl.accessChan <- accessConfig
			}
		}
	}()
}

func (ctrl *GWController) Stop() {
	close(ctrl.stopC)
}
