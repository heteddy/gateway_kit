// @Author : detaohe
// @File   : access.go
// @Description:
// @Date   : 2022/9/11 16:31

package gateway

import (
	"gateway_kit/config"
	"gateway_kit/dao"
	"gateway_kit/util"
	"go.uber.org/zap"
	"sync"
)

const (
	ACCESS_CONTROL_GATEWAY = iota
	ACCESS_CONTROL_SERVICE
)

type AccessConfig struct {
	Name      string
	EventType int
	Category  int
	BlockIP   []string
	AllowIP   []string
}

var Access *AccessController
var onceAccess sync.Once

type AccessController struct {
	mutex      sync.RWMutex
	svcBlockIP map[string][]string
	svcAllowIP map[string][]string
	gwBlockIP  []string
	gwAllowIP  []string

	accessChan chan *AccessConfig
	stopC      chan struct{}
}

func NewAccessController() *AccessController {
	onceAccess.Do(func() {
		Access = &AccessController{
			mutex:      sync.RWMutex{},
			svcBlockIP: make(map[string][]string),
			svcAllowIP: make(map[string][]string),
			gwBlockIP:  nil,
			gwAllowIP:  nil,
			accessChan: make(chan *AccessConfig),
			stopC:      make(chan struct{}),
		}
		Access.Start()
	})
	return Access

}

func (ac *AccessController) update(c *AccessConfig) {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()
	switch c.Category {
	case ACCESS_CONTROL_GATEWAY: // gateway 只可能有一个
		// gateway不支持删除
		ac.gwAllowIP = c.AllowIP
		ac.gwBlockIP = c.BlockIP
	case ACCESS_CONTROL_SERVICE:
		if c.EventType == dao.EventDelete {
			delete(ac.svcAllowIP, c.Name)
			delete(ac.svcBlockIP, c.Name)
		} else {
			ac.svcAllowIP[c.Name] = c.AllowIP
			ac.svcBlockIP[c.Name] = c.BlockIP
		}

	default:
		config.Logger.Warn("receiving error category", zap.Int("category", c.Category))
	}
}

func (ac *AccessController) runLoop() {
loop:
	for {
		select {
		case c, ok := <-ac.accessChan:
			if !ok {
				break loop
			}
			ac.update(c)
		case <-ac.stopC:
			break loop
		}
	}
}

func (ac *AccessController) Start() {
	go ac.runLoop()
}

func (ac *AccessController) In() chan<- *AccessConfig {
	return ac.accessChan
}

func (ac *AccessController) Stop() {
	close(ac.stopC)
	//close(ac.accessChan) // todo 这里有点危险，应该由写入的关闭
}

func (ac *AccessController) IsAllowed(svc, ip string) bool {
	ac.mutex.RLock()
	defer ac.mutex.RUnlock()

	_block := util.IPSlice(ac.gwBlockIP)
	if _block.Has(ip) {
		return false
	}
	_block = ac.gwAllowIP
	if _block.Has(ip) {
		return true
	}
	_block = ac.svcBlockIP[svc]
	if _block.Has(ip) {
		return false
	}
	_block = ac.svcAllowIP[svc]
	if _block.Has(ip) {
		return true
	}
	return true
}
