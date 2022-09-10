// @Author : detaohe
// @File   : random.go
// @Description:
// @Date   : 2022/9/3 18:55

package lb

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type randomLB struct {
	mutex        sync.RWMutex
	serviceAddrs map[string][]string
}

func NewRandomLB() LoadBalancer {
	return &randomLB{
		mutex:        sync.RWMutex{},
		serviceAddrs: make(map[string][]string),
	}
}
func (lb *randomLB) Update(svc string, addrs []string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	if _, existed := lb.serviceAddrs[svc]; existed {
		lb.serviceAddrs[svc] = addrs
	}
}
func (lb *randomLB) Next(svc string) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Int63()
	lb.mutex.RLock()
	addrs := lb.serviceAddrs[svc]
	lb.mutex.RUnlock()
	if len(addrs) > 0 {
		idx := n % int64(len(addrs))
		return addrs[idx], nil
	} else {
		return "", errors.New("no service available")
	}
}
