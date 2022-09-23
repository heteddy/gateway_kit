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
	mutex    sync.RWMutex
	nodeMap  map[string][]*Node
	svcIndex map[string]int64
}

func NewRoundRobin() LoadBalancer {
	return &roundRobin{
		mutex:    sync.RWMutex{},
		nodeMap:  make(map[string][]*Node),
		svcIndex: make(map[string]int64),
	}
}

func (lb *roundRobin) Next(svc string) (string, error) {
	lb.mutex.RLock()
	count := lb.svcIndex[svc]
	lb.mutex.RUnlock()
	count++ //这里有可能在并发情况下，导致count实际不准, 这里不需要写锁，并发量优先
	nodes := lb.nodeMap[svc]
	if len(nodes) > 0 {
		idx := count % int64(len(nodes))
		lb.svcIndex[svc] = count //这里有可能在并发情况下，导致count实际不准
		return nodes[idx].Addr, nil
	} else {
		return "", errors.New("service " + svc + " not found")
	}
}

func (lb *roundRobin) UpdateNode(nodes *Node) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	//for _, n := range nodes {
	//	if nodeList, existed := lb.nodeMap[n.Svc]; existed {
	//		lb.nodeMap[n.Svc] = append(nodeList, n)
	//	} else {
	//		newNodeList := make([]*Node, 0, 1)
	//		newNodeList = append(newNodeList, n)
	//		lb.nodeMap[n.Svc] = newNodeList
	//		lb.svcIndex[n.Svc] = 0
	//	}
	//}
}
