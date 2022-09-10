// @Author : detaohe
// @File   : roundrobin.go
// @Description:
// @Date   : 2022/9/3 17:24

package lb

import (
	"errors"
	"sync"
)

type roundRobin struct {
	mutex        sync.RWMutex
	serviceAddrs map[string][]string
	serviceIndex map[string]int64
}

func (rb *roundRobin) Next(svc string) (string, error) {
	rb.mutex.RLock()
	count := rb.serviceIndex[svc]
	rb.mutex.RUnlock()
	count++ //这里有可能在并发情况下，导致count实际不准, 这里不需要写锁，并发量优先
	addrs := rb.serviceAddrs[svc]
	if len(addrs) > 0 {
		idx := count % int64(len(addrs))
		rb.serviceIndex[svc] = count //这里有可能在并发情况下，导致count实际不准
		return addrs[idx], nil
	} else {
		return "", errors.New("no service available")
	}
}

func (rb *roundRobin) Update(svc string, addrs []string) {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()
	rb.serviceAddrs[svc] = addrs
	if _, existed := rb.serviceIndex[svc]; !existed {
		rb.serviceIndex[svc] = 0
	}
}

func NewRoundRobin() LoadBalancer {
	return &roundRobin{
		mutex:        sync.RWMutex{},
		serviceAddrs: make(map[string][]string),
		serviceIndex: make(map[string]int64),
	}
}
