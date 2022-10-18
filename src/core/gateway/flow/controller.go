// @Author : detaohe
// @File   : controller.go
// @Description:
// @Date   : 2022/10/15 16:30

package flow

import (
	"context"
	"fmt"
	"gateway_kit/config"
	"gateway_kit/dao"
	"gateway_kit/util"
	"go.uber.org/zap"
	"strings"
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
	config.Logger.Info("add gw flow", zap.String("key", ctrl.Key()), zap.Int64("value", 1))
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
	config.Logger.Info("clean gw flow", zap.String("key", ctrl.Key()), zap.Int64("value", value))
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

func (ctrl *SvcFlowCtrl) EncodeKey() string {
	return fmt.Sprintf("%s:%s", ctrl.keyPrefix, ctrl.name)
}
func (ctrl *SvcFlowCtrl) RequestPrefix() string {
	return fmt.Sprintf("%s:", ctrl.EncodeKey()) //note 这里request不能等于service 否则会把service的找出来
}
func (ctrl *SvcFlowCtrl) EncodeRequestKey(uri, method string) string {
	return fmt.Sprintf("%s%s:%s", ctrl.RequestPrefix(), uri, method)
}
func (ctrl *SvcFlowCtrl) DecodeRequestKey(key string) (svc, uri, method string) {
	ss := strings.Split(key, ":")
	if len(ss) == 4 {
		svc, uri, method = ss[1], ss[2], ss[3]
	} else {
		config.Logger.Warn("request key length", zap.Int("length", len(ss)), zap.String("key", key))
	}
	return
}

func (ctrl *SvcFlowCtrl) Call(uri, method string) {
	ctrl.storage.In() <- &IncreaseSvcFlowCmd{
		svcKey:   ctrl.EncodeKey(),
		svcValue: 1,
		reqKey:   ctrl.EncodeRequestKey(uri, method),
		reqValue: 1,
	}
}

func (ctrl *SvcFlowCtrl) Load(svcC, reqC chan map[string]int64) {
	// 查询service的值
	// 根据request prefix 查询所有的request 值
	ctrl.storage.In() <- &LoadSvcFlowCmd{
		svcKey:    ctrl.EncodeKey(),
		reqPrefix: ctrl.RequestPrefix() + "*",
		svcC:      svcC,
		reqC:      reqC,
	}
}

func (ctrl *SvcFlowCtrl) Clean(svc, request map[string]int64) {
	// 通常写到数据库之后，要求清掉redis的数据
	config.Logger.Info("clean svc flow", zap.String("key", ctrl.EncodeKey()), zap.Int64("value", svc[ctrl.EncodeKey()]))
	ctrl.storage.In() <- &DecreaseSvcFlowCmd{
		svcKey:      ctrl.EncodeKey(),
		svcValue:    svc[ctrl.EncodeKey()],
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
	svcHourDao *dao.ServiceHourDao
	reqHourDao *dao.ReqHourDao
}

func NewFlowCollector() *FlowCollector {
	flowOnce.Do(func() {
		gw := config.All.Name
		storage := NewFlowStorage()
		flowCollector = &FlowCollector{
			mutex:  sync.RWMutex{},
			gwName: gw,
			gwFlow: &GwFlowCtrl{
				name:    gw,
				key:     "flow_" + gw,
				storage: storage,
			},
			svcFlows:   make(map[string]*SvcFlowCtrl),
			storage:    storage,
			stopC:      make(chan struct{}),
			TickerSvc:  util.NewTickerSvc("FlowCollector", time.Second*10, false),
			svcHourDao: dao.NewServiceHourDao(),
			reqHourDao: dao.NewReqHourDao(),
		}
	})
	return flowCollector
}

func (collector *FlowCollector) Call(svc, uri, method string) {
	collector.gwFlow.Call()
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
		config.Logger.Info("定时刷新redis flow到mongo")
		// 定时任务刷新redis数据到mongo， 同时删除redis中的内容
		// note 这个任务不能多进程并发处理，因为redis的decrease可能导致重复
		// 1. 获取所有的gw和service的访问数量
		gwC := make(chan map[string]int64)
		defer close(gwC)
		collector.gwFlow.Load(gwC)
		gwInfo := <-gwC

		gwCount := gwInfo[collector.gwFlow.Key()]
		config.Logger.Info("gw flow", zap.Int64("gwCount", gwCount), zap.String("gwname", collector.gwName))

		now := time.Now()

		collector.svcHourDao.Update(context.Background(), &dao.ServiceHourEntity{
			Category: "gateway",
			Name:     collector.gwName, // http服务名词
			Count:    gwCount,
			Hour:     now.Hour(), // http
			Date:     now.Format("2006-01-02"),
		})

		collector.gwFlow.Clean(gwCount) // 清除gateway调用次数

		collector.mutex.Lock()
		for name, svcFC := range collector.svcFlows {
			svcC := make(chan map[string]int64)
			reqC := make(chan map[string]int64)
			svcFC.Load(svcC, reqC)

			svcKv := <-svcC
			for _, v := range svcKv {
				// 这里只可能有一个
				if v == 0 {
					continue
				}
				collector.svcHourDao.Update(context.Background(), &dao.ServiceHourEntity{
					Category: "service",
					Name:     name, // http服务名词
					Count:    v,
					Hour:     now.Hour(), // http
					Date:     now.Format("2006-01-02"),
				})
			}
			fmt.Println("svcKv:", svcKv)
			reqKv := <-reqC
			for k, v := range reqKv {
				if v == 0 {
					continue
				}
				svc, path, method := svcFC.DecodeRequestKey(k)
				collector.reqHourDao.Update(context.Background(), &dao.ReqHourEntity{
					Service: svc,
					Path:    path, // http服务名词
					Method:  method,
					Count:   v,
					Hour:    now.Hour(), // http
					Date:    now.Format("2006-01-02"),
				})
			}
			svcFC.Clean(svcKv, reqKv) // 清除服务调用次数
			fmt.Println("reqKv:", reqKv)
			close(svcC)
			close(reqC)
		}
		collector.mutex.Unlock()
	}
}
