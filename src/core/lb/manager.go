// @Author : detaohe
// @File   : manager.go
// @Description:
// @Date   : 2022/9/11 20:33

package lb

import (
	"gateway_kit/config"
	"go.uber.org/zap"
	"sync"
)

const (
	LbRandom     = "random"
	LbWeight     = "weight"
	LbRoundRobin = "round_robin"
)

type LoadBalanceFactory struct {
}

var lbManager *LoadBalanceMgr
var lbOnce sync.Once

func NewLBManager() *LoadBalanceMgr {
	lbOnce.Do(func() {
		lbManager = &LoadBalanceMgr{
			balancerMap: make(map[string]LoadBalancer),
			addrChan:    make(chan *Node),
			stopC:       make(chan struct{}),
		}
		lbManager.init()
	})
	return lbManager
}

func (f *LoadBalanceFactory) Create(_type string) LoadBalancer {
	switch _type {
	case LbRandom:
		return NewRandomLB()
	case LbWeight:
		return NewWeightedRoundRobinLB()
	case LbRoundRobin:
		return NewRoundRobin()
	default:
		return nil
	}
}

type LoadBalanceMgr struct {
	balancerMap map[string]LoadBalancer
	addrChan    chan *Node
	stopC       chan struct{}
	//mutex       sync.RWMutex
}

//func NewLoadBalanceMgr() *LoadBalanceMgr {
//	lbm := &LoadBalanceMgr{
//		balancerMap: make(map[string]LoadBalancer),
//		addrChan:    make(chan *Node),
//		stopC:       make(chan struct{}),
//		//mutex:       sync.RWMutex{},
//	}
//	lbm.init()
//	return lbm
//}

func (m *LoadBalanceMgr) init() {
	loadTypes := []string{
		LbRandom, LbWeight, LbRoundRobin,
	}
	factory := LoadBalanceFactory{}
	for _, t := range loadTypes {
		m.balancerMap[t] = factory.Create(t)
	}
}

func (m *LoadBalanceMgr) In() chan<- *Node {
	return m.addrChan
}

func (m *LoadBalanceMgr) Get(lbType string) LoadBalancer {
	if _lb, existed := m.balancerMap[lbType]; existed {
		return _lb
	}
	return nil
}

func (m *LoadBalanceMgr) update(node *Node) {
	config.Logger.Info("update LoadBalanceMgr", zap.Any("Node", node))
	for k, v := range m.balancerMap {
		config.Logger.Info("update lb", zap.String("lb_type", k))
		v.UpdateNode(node)
	}
}

func (m *LoadBalanceMgr) runLoop() {
loop:
	for {
		select {
		case node, ok := <-m.addrChan:
			if !ok {
				config.Logger.Warn("LoadBalanceMgr exit")
				break loop
			}
			m.update(node)
		case <-m.stopC:
			config.Logger.Warn("LoadBalanceMgr exit")
			break loop
		}
	}
}

func (m *LoadBalanceMgr) Start() {
	go m.runLoop()
}

func (m *LoadBalanceMgr) Stop() {
	close(m.stopC)
	//close(m.addrChan) // todo 这里有点危险，应该由写入的关闭
}
