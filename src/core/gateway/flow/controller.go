// @Author : detaohe
// @File   : controller.go
// @Description:
// @Date   : 2022/10/15 16:30

package flow

import (
	"fmt"
	"gateway_kit/util"
	"sync"
	"time"
)

var flowCollector *FlowCollector
var flowOnce sync.Once

type GwFlowCtrl struct {
	name    string
	key     string
	storage *FlowStorage
}

func (ctrl *GwFlowCtrl) Key() string {
	return fmt.Sprintf("%s", ctrl.key)
}

func (ctrl *GwFlowCtrl) Call() {
	ctrl.storage.In() <- &IncreaseGwFlowCmd{
		gwKey:   ctrl.Key(),
		gwValue: 1,
	}
}

func (ctrl *GwFlowCtrl) Load(gwC chan map[string]int64) {
	// 查询service的值
	// 根据request prefix 查询所有的request 值
	ctrl.storage.In() <- &LoadGwFlowCmd{
		key: ctrl.Key(),
		gwC: gwC,
	}
}

func (ctrl *GwFlowCtrl) Clean(value int64) {
	// 通常写到数据库之后，要求清掉redis的数据
	ctrl.storage.In() <- &DecreaseGwFlowCmd{
		gwKey:   ctrl.Key(),
		gwValue: value,
	}
}

// SvcFlowCtrl 刷新到redis
type SvcFlowCtrl struct {
	name      string
	keyPrefix string
	storage   *FlowStorage
}

func (ctrl *SvcFlowCtrl) Key() string {
	return fmt.Sprintf("%s:%s", ctrl.keyPrefix, ctrl.name)
}
func (ctrl *SvcFlowCtrl) RequestPrefix() string {
	return fmt.Sprintf("%s:req", ctrl.Key())
}
func (ctrl *SvcFlowCtrl) RequestKey(uri, method string) string {
	return fmt.Sprintf("%s:%s:%s", ctrl.RequestPrefix(), uri, method)
}

func (ctrl *SvcFlowCtrl) Call(uri, method string) {
	ctrl.storage.In() <- &IncreaseSvcFlowCmd{
		svcKey:   ctrl.Key(),
		svcValue: 1,
		reqKey:   ctrl.RequestKey(uri, method),
		reqValue: 1,
	}
}

func (ctrl *SvcFlowCtrl) Load(svcC, reqC chan map[string]int64) {
	// 查询service的值
	// 根据request prefix 查询所有的request 值
	ctrl.storage.In() <- &LoadSvcFlowCmd{
		svcKey:    ctrl.Key(),
		reqPrefix: ctrl.RequestPrefix(),
		svcC:      svcC,
		reqC:      reqC,
	}
}

func (ctrl *SvcFlowCtrl) Clean(svc, request map[string]int64) {
	// 通常写到数据库之后，要求清掉redis的数据
	ctrl.storage.In() <- &DecreaseSvcFlowCmd{
		svcKey:      ctrl.Key(),
		svcValue:    svc[ctrl.Key()],
		requestFlow: request,
	}
}

// FlowCollector 定时从redis刷新到mongodb
type FlowCollector struct {
	mutex    sync.RWMutex
	gwName   string
	gwFlow   *GwFlowCtrl
	svcFlows map[string]*SvcFlowCtrl // service
	storage  *FlowStorage
	stopC    chan struct{}
	*util.TickerSvc
}

func NewFlowCollector(gw string) *FlowCollector {
	flowOnce.Do(func() {
		storage := NewFlowStorage()
		flowCollector = &FlowCollector{
			mutex:  sync.RWMutex{},
			gwName: gw,
			gwFlow: &GwFlowCtrl{
				name:    gw,
				key:     "flow_" + gw,
				storage: storage,
			},
			svcFlows:  make(map[string]*SvcFlowCtrl),
			storage:   storage,
			stopC:     make(chan struct{}),
			TickerSvc: util.NewTickerSvc("FlowCollector", time.Minute*1, false),
		}
	})
	return flowCollector
}

func (collector *FlowCollector) Call(svc, uri, method string) {
	collector.mutex.Lock()
	defer collector.mutex.Unlock()
	if svcFlow, existed := collector.svcFlows[svc]; existed {
		svcFlow.Call(uri, method)
	} else {
		svcFlow := &SvcFlowCtrl{
			name:      svc,
			keyPrefix: "flow_" + collector.gwName,
			storage:   collector.storage,
		}
		collector.svcFlows[svc] = svcFlow
		svcFlow.Call(uri, method)
	}
}

func (collector *FlowCollector) Stop() {
	collector.TickerSvc.Stop()
	close(collector.stopC)
	close(collector.storage.storageChan)
}
func (collector *FlowCollector) Start() {
	collector.storage.Start()
	collector.TickerSvc.Start(collector.endpoint())
}

func (collector *FlowCollector) endpoint() util.SvcEndpoint {
	return func() {
		// 定时任务刷新redis数据到mongo， 同时删除redis中的内容
		// note 这个任务不能多进程并发处理，因为redis的decrease可能导致重复
	}
}
