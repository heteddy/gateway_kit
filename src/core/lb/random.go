// @Author : detaohe
// @File   : random.go
// @Description:
// @Date   : 2022/9/3 18:55

package lb

import (
	"errors"
	"fmt"
	"gateway_kit/config"
	"gateway_kit/dao"
	"go.uber.org/zap"
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
func (lb *randomLB) Log() {
	fmt.Printf("\nstarting log routing table\n")
	for k, _nodes := range lb.serviceAddrs {
		for _, v := range _nodes {
			fmt.Printf("name=%20s, node=%v\n", k, v)
		}
	}
	fmt.Printf("\n")
}
func (lb *randomLB) UpdateNode(node *Node) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	//lb.serviceAddrs = make(map[string][]*Node)
	// service name 存在
	config.Logger.Info("updating node", zap.Any("node", node))
	if _nodes, existed := lb.serviceAddrs[node.Svc]; existed {
		switch node.EventType {
		case dao.EventDelete:
		loop:
			for idx, n := range _nodes {
				if n.IsSameNode(*node) {
					_addrList := append(_nodes[0:idx], _nodes[idx+1:]...)
					lb.serviceAddrs[node.Svc] = _addrList
					break loop
				}
			}
		case dao.EventUpdate:
			for _, n := range _nodes {
				if n.IsSameNode(*node) {
					n.Update(*node)
				}
			}
		case dao.EventCreate:
			nodes := append(_nodes, node)
			lb.serviceAddrs[node.Svc] = nodes
		default:
		}
	} else { // 没有找到，添加到
		if node.EventType != dao.EventDelete {
			newNodes := make([]*Node, 1, 1)
			newNodes[0] = node
			lb.serviceAddrs[node.Svc] = newNodes
		}
	}
	lb.Log()
}
func (lb *randomLB) Next(svc string) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Int63()
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()
	if addrs, existed := lb.serviceAddrs[svc]; existed {
		if len(addrs) > 0 {
			idx := n % int64(len(addrs))
			return addrs[idx].Addr, nil
		} else {
			return "", errors.New("no service available")
		}
	} else {
		return "", errors.New("service not existed")
	}

}
