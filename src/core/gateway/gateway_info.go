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
	gwChan     chan *dao.GatewayEntity
	accessChan chan []*AccessConfig
}

func NewGwConfig(accessChan chan []*AccessConfig) *GWConfig {
	onceGwConfig.Do(func() {
		GwConfigure = &GWConfig{
			mutex:      sync.RWMutex{},
			stopC:      make(chan struct{}),
			gwChan:     make(chan *dao.GatewayEntity),
			accessChan: accessChan,
		}
		GwConfigure.Start()
	})
	return GwConfigure
}

func (configure *GWConfig) In() chan *dao.GatewayEntity {
	return configure.gwChan
}
func (configure *GWConfig) Start() {
	go func() {
	loop:
		for {
			select {
			case <-configure.stopC:
				break loop
			case entity, ok := <-configure.gwChan:
				if !ok {
					break loop
				}
				accessConfig := &AccessConfig{
					BlockIP:  entity.BlockList,
					AllowIP:  entity.AllowList,
					Name:     entity.Name,
					Category: ACCESS_CONTROL_GATEWAY,
				}
				configure.accessChan <- []*AccessConfig{accessConfig}
			}
		}
	}()
}

func (configure *GWConfig) Stop() {
	close(configure.stopC)
}
