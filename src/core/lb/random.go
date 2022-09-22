// @Author : detaohe
// @File   : random.go
// @Description:
// @Date   : 2022/9/3 18:55

package lb

import (
	"errors"
	"gateway_kit/dao"
	"math/rand"
	"sync"
	"time"
)

type randomLB struct {
	mutex        sync.RWMutex
	serviceAddrs map[string][]*Node
}

func NewRandomLB() LoadBalancer {
	return &randomLB{
		mutex:        sync.RWMutex{},
		serviceAddrs: make(map[string][]*Node),
	}
}
func (lb *randomLB) UpdateNode(node *Node) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	//lb.serviceAddrs = make(map[string][]*Node)
	// service name 存在
	if _nodes, existed := lb.serviceAddrs[node.Svc]; existed {
		if node.EventType == dao.EventDelete {
			for
		}
		lb.serviceAddrs[node.Svc] = append(_nodes, node)
	} else {
		newNodes := make([]*Node, 1, 1)
		newNodes[0] = node
		lb.serviceAddrs[node.Svc] = newNodes
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
		return addrs[idx].Addr, nil
	} else {
		return "", errors.New("no service available")
	}
}
