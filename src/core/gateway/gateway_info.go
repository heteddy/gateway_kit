// @Author : detaohe
// @File   : gateway_info.go
// @Description:
// @Date   : 2022/9/13 15:42

package gateway

import (
	"gateway_kit/dao"
	"sync"
)

var onceGwConfig sync.Once
var GwConfigure *GWConfig

type GWConfig struct {
	mutex      sync.RWMutex
	stopC      chan struct{}
	gwChan     chan *dao.GwEvent
	accessChan chan *AccessConfig
}

func NewGwConfig(accessChan chan *AccessConfig) *GWConfig {
	onceGwConfig.Do(func() {
		GwConfigure = &GWConfig{
			mutex:      sync.RWMutex{},
			stopC:      make(chan struct{}),
			gwChan:     make(chan *dao.GwEvent),
			accessChan: accessChan,
		}
		GwConfigure.Start()
	})
	return GwConfigure
}

func (configure *GWConfig) In() chan *dao.GwEvent {
	return configure.gwChan
}
func (configure *GWConfig) Start() {
	go func() {
	loop:
		for {
			select {
			case <-configure.stopC:
				break loop
			case event, ok := <-configure.gwChan:
				if !ok {
					break loop
				}
				for _, entity := range event.Entities {
					accessConfig := &AccessConfig{
						EventType: event.EventType,
						BlockIP:   entity.BlockList,
						AllowIP:   entity.AllowList,
						Name:      entity.Name,
						Category:  ACCESS_CONTROL_GATEWAY,
					}
					configure.accessChan <- accessConfig
				}
			}
		}
	}()
}

func (configure *GWConfig) Stop() {
	close(configure.stopC)
}
