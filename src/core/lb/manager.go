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
