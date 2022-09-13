// @Author : detaohe
// @File   : transport
// @Description:
// @Date   : 2022/9/9 14:31

package gateway

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var TransportPoolGen *TransportPool
var onceTransport sync.Once

type TransportPool struct {
	transportMap map[string]*http.Transport
	mutex        sync.RWMutex
}

func NewTransportPool() *TransportPool {
	onceTransport.Do(func() {
		TransportPoolGen = &TransportPool{
			transportMap: make(map[string]*http.Transport),
			mutex:        sync.RWMutex{},
		}
	})
	return TransportPoolGen
}

func (tp *TransportPool) Get(svcName string) *http.Transport {
	tp.mutex.RLock()
	v, existed := tp.transportMap[svcName]
	tp.mutex.RUnlock()
	if existed {
		return v
	} else {
		tp.mutex.Lock()
		defer tp.mutex.Unlock()
		_newTp := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(3) * time.Second,
				KeepAlive: 60 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          50,
			IdleConnTimeout:       time.Duration(60) * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: time.Duration(3) * time.Second,
		}
		tp.transportMap[svcName] = _newTp
		return _newTp
	}
}
