// @Author : detaohe
// @File   : protocol
// @Description:
// @Date   : 2022/9/22 15:33

package gateway

import (
	"gateway_kit/dao"
	"sync"
)

var ProtoTrans *ProtocolTrans
var onceProto sync.Once

type ProtocolSupported struct {
	EventType   int
	Svc         string
	IsHttps     bool `json:"need_https"  description:"type=支持https 1=支持"`
	IsWebsocket bool `json:"need_websocket" description:"启用websocket 1=启用"`
}

type ProtocolTrans struct {
	mutex        sync.RWMutex
	svcProtocols map[string]*ProtocolSupported

	protoChan chan *ProtocolSupported
	stopC     chan struct{}
}

func NewProtocolTrans() *ProtocolTrans {
	onceProto.Do(func() {
		ProtoTrans = &ProtocolTrans{
			mutex:        sync.RWMutex{},
			svcProtocols: make(map[string]*ProtocolSupported),
			protoChan:    make(chan *ProtocolSupported),
			stopC:        make(chan struct{}),
		}
	})
	return ProtoTrans
}

func (proto *ProtocolTrans) runLoop() {
loop:
	for {
		select {
		case <-proto.stopC:
			break loop
		case protocol, ok := <-proto.protoChan:
			if !ok {
				break loop
			}
			proto.update(protocol)
		}
	}
}
func (proto *ProtocolTrans) Start() {
	go proto.runLoop()
}

func (proto *ProtocolTrans) In() chan<- *ProtocolSupported {
	return proto.protoChan
}

func (proto *ProtocolTrans) Stop() {
	close(proto.stopC)
}

func (proto *ProtocolTrans) update(protocol *ProtocolSupported) {
	proto.mutex.Lock()
	defer proto.mutex.Unlock()

	if protocol.EventType == dao.EventDelete {
		// 删除
		delete(proto.svcProtocols, protocol.Svc)
	} else {
		proto.svcProtocols[protocol.Svc] = protocol
	}
}

func (proto *ProtocolTrans) IsHttps(name string) bool {
	proto.mutex.RLock()
	defer proto.mutex.RUnlock()
	if p, existed := proto.svcProtocols[name]; existed {
		return p.IsHttps
	}
	return false
}

func (proto *ProtocolTrans) IsWebsocket(name string) bool {
	proto.mutex.RLock()
	defer proto.mutex.RUnlock()
	if p, existed := proto.svcProtocols[name]; existed {
		return p.IsWebsocket
	}
	return false
}
