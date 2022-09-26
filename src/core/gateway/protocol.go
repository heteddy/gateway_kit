// @Author : detaohe
// @File   : protocol
// @Description:
// @Date   : 2022/9/22 15:33

package gateway

import (
	"gateway_kit/config"
	"gateway_kit/dao"
	"go.uber.org/zap"
	"sync"
)

var ProtoTrans *ProtocolTransCtrl
var onceProto sync.Once

type ProtocolSupportedEvt struct {
	EventType   int
	Svc         string
	IsHttps     bool `json:"is_https"  description:"type=支持https 1=支持"`
	IsWebsocket bool `json:"is_websocket" description:"启用websocket 1=启用"`
}

type ProtocolTransCtrl struct {
	mutex        sync.RWMutex
	svcProtocols map[string]*ProtocolSupportedEvt

	protoChan chan *ProtocolSupportedEvt
	stopC     chan struct{}
}

func NewProtocolTransCtrl() *ProtocolTransCtrl {
	onceProto.Do(func() {
		ProtoTrans = &ProtocolTransCtrl{
			mutex:        sync.RWMutex{},
			svcProtocols: make(map[string]*ProtocolSupportedEvt),
			protoChan:    make(chan *ProtocolSupportedEvt),
			stopC:        make(chan struct{}),
		}
	})
	return ProtoTrans
}

func (proto *ProtocolTransCtrl) runLoop() {
loop:
	for {
		select {
		case <-proto.stopC:
			config.Logger.Warn("ProtocolTransCtrl exit")
			break loop
		case protocol, ok := <-proto.protoChan:
			if !ok {
				config.Logger.Warn("ProtocolTransCtrl exit")
				break loop
			}
			proto.update(protocol)
		}
	}
}
func (proto *ProtocolTransCtrl) Start() {
	go proto.runLoop()
}

func (proto *ProtocolTransCtrl) In() chan<- *ProtocolSupportedEvt {
	return proto.protoChan
}

func (proto *ProtocolTransCtrl) Stop() {
	close(proto.stopC)
}

func (proto *ProtocolTransCtrl) update(protocol *ProtocolSupportedEvt) {
	proto.mutex.Lock()
	defer proto.mutex.Unlock()
	config.Logger.Info("update ProtocolTransCtrl", zap.Any("ProtocolSupportedEvt", protocol))
	if protocol.EventType == dao.EventDelete {
		// 删除
		delete(proto.svcProtocols, protocol.Svc)
	} else {
		proto.svcProtocols[protocol.Svc] = protocol
	}
}

func (proto *ProtocolTransCtrl) IsHttps(name string) bool {
	proto.mutex.RLock()
	defer proto.mutex.RUnlock()
	if p, existed := proto.svcProtocols[name]; existed {
		return p.IsHttps
	}
	return false
}

func (proto *ProtocolTransCtrl) IsWebsocket(name string) bool {
	proto.mutex.RLock()
	defer proto.mutex.RUnlock()
	if p, existed := proto.svcProtocols[name]; existed {
		return p.IsWebsocket
	}
	return false
}
