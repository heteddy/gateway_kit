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
	Lb_Random     = "random"
	Lb_Weight     = "weight"
	Lb_RoundRobin = "round_robin"
)

type LoadBalanceFactory struct {
}

var Manager *LoadBalanceMgr

//var factory *LoadBalanceFactory
var once sync.Once

func NewManager() *LoadBalanceMgr {
	once.Do(func() {

		Manager = &LoadBalanceMgr{
			balancerMap: make(map[string]LoadBalancer),
		}
		Manager.init()
	})
	return Manager
}

func (f *LoadBalanceFactory) Create(_type string) LoadBalancer {
	switch _type {
	case Lb_Random:
		return NewRandomLB()
	case Lb_Weight:
		return NewWeightedRoundRobinLB()
	case Lb_RoundRobin:
		return NewRoundRobin()
	default:
		return nil
	}
}

type LoadBalanceMgr struct {
	balancerMap map[string]LoadBalancer
	addrChan    chan []*Node
	stopC       chan struct{}
	mutex       sync.RWMutex
}

func NewLoadBalanceMgr() *LoadBalanceMgr {
	lbm := &LoadBalanceMgr{
		balancerMap: make(map[string]LoadBalancer),
		addrChan:    make(chan []*Node),
		stopC:       make(chan struct{}),
		mutex:       sync.RWMutex{},
	}
	lbm.init()
	return lbm
}

func (m *LoadBalanceMgr) init() {
	loadTypes := []string{
		Lb_Random, Lb_Weight, Lb_RoundRobin,
	}
	factory := LoadBalanceFactory{}
	for _, t := range loadTypes {
		m.balancerMap[t] = factory.Create(t)
	}
}

func (m *LoadBalanceMgr) Update(nodes []*Node) {
	for k, v := range m.balancerMap {
		config.Logger.Info("update lb", zap.String("lb_type", k))
		v.UpdateNodes(nodes)
	}
}

func (m *LoadBalanceMgr) runLoop() {
loop:
	for {
		select {
		case addrs, ok := <-m.addrChan:
			if !ok {
				break loop
			}
			m.Update(addrs)
		case <-m.stopC:
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
