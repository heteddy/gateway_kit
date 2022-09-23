// @Author : detaohe
// @File   : svr_proxy.go
// @Description:
// @Date   : 2022/9/12 09:28

package discovery

import (
	"context"
	"gateway_kit/dao"
	"gateway_kit/util"
	"time"
)

// 暂时通过admin的service 获取服务，后期可以根据数据库和

type ServiceDiscovery interface {
}

type PollingDiscovery struct {
	svcDao *dao.HttpSvcDao // 定时获取所有的服务
	gwDao  *dao.GatewayDao
	*util.TickerSvc
	svcChan   chan<- *dao.SvcEvent
	gwChan    chan *dao.GwEvent
	gwInfo    *dao.GatewayEntity
	pollingAt *time.Time
	//mutex       sync.RWMutex
	//entities     []*dao.HttpSvcEntity
	//svcEntityMap map[string][]*dao.HttpSvcEntity // 以后放到redis，目前
}

func NewPollingDiscovery(svcChan chan *dao.SvcEvent, gwChan chan *dao.GwEvent) *PollingDiscovery {
	proxy := &PollingDiscovery{
		svcDao:    dao.NewHttpSvcDao(),
		gwDao:     dao.NewGatewayDao(),
		TickerSvc: util.NewTickerSvc("service-discovery", time.Minute, false),
		svcChan:   svcChan,
		gwChan:    gwChan,
		pollingAt: nil,
		//entities:     make([]*dao.HttpSvcEntity, 0, 1),
		//svcEntityMap: make(map[string][]*dao.HttpSvcEntity),
		//mutex:       sync.RWMutex{},
	}
	proxy.Start(proxy.endpoint())
	return proxy
}

func (pd *PollingDiscovery) loadServices() ([]*dao.HttpSvcEntity, error) {
	if pd.pollingAt == nil {
		return pd.svcDao.All(context.Background())
	} else {
		return pd.svcDao.ChangeFrom(context.Background(), *pd.pollingAt)
	}
}

func (pd *PollingDiscovery) loadGateway() ([]*dao.GatewayEntity, error) {
	if pd.pollingAt == nil {
		return pd.gwDao.All(context.Background())
	} else {
		return pd.gwDao.ChangeFrom(context.Background(), *pd.pollingAt)
	}
}

func (pd *PollingDiscovery) endpoint() util.SvcEndpoint {
	return func() {
		//entities, err := pd.svcDao.All(context.Background())
		//if err != nil {
		//	config.Logger.Error("load service entities error", zap.Error(err))
		//	return
		//}
		//pd.mutex.Lock()
		//defer pd.mutex.Unlock()
		//pd.svcEntityMap = make(map[string][]*dao.HttpSvcEntity)
		//for _, e := range entities {
		//	if v, existed := pd.svcEntityMap[e.Name]; existed {
		//		pd.svcEntityMap[e.Name] = append(v, e)
		//	} else {
		//		_newList := make([]*dao.HttpSvcEntity, 0, 1)
		//		_newList = append(_newList, e)
		//		pd.svcEntityMap[e.Name] = _newList
		//	}
		//}
		//pd.entities = entities
		//if len(pd.entities) > 0 {
		//	pd.svcChan <- pd.entities
		//}
		//
		//// 读取gateway
		//gatewayInfos, err2 := pd.gwDao.GetByName(context.Background(), config.All.Name)
		//if err2 != nil {
		//	config.Logger.Error("load gateway entities error", zap.Error(err2))
		//	return
		//} else {
		//	if len(gatewayInfos) > 0 {
		//		pd.gwInfo = gatewayInfos[0]
		//		pd.gwChan <- pd.gwInfo
		//	}
		//}
		if svcs, err1 := pd.loadServices(); err1 != nil {

		} else {
			delEvents := &dao.SvcEvent{
				EventType: dao.EventDelete,
				Entities:  make([]*dao.HttpSvcEntity, 0),
			}
			updateEvents := &dao.SvcEvent{
				EventType: dao.EventUpdate,
				Entities:  make([]*dao.HttpSvcEntity, 0),
			}
			for _, svc := range svcs {
				if svc.DeletedAt > 0 {
					delEvents.Entities = append(delEvents.Entities, svc)
				} else {
					updateEvents.Entities = append(updateEvents.Entities, svc)
				}
			}
			if !delEvents.Empty() {
				pd.svcChan <- delEvents
			}
			if !updateEvents.Empty() {
				pd.svcChan <- updateEvents
			}
		}

		if gws, err2 := pd.loadGateway(); err2 != nil {

		} else {
			delEvent := &dao.GwEvent{
				EventType: dao.EventDelete,
				Entities:  make([]*dao.GatewayEntity, 0),
			}
			updateEvent := &dao.GwEvent{
				EventType: dao.EventUpdate,
				Entities:  make([]*dao.GatewayEntity, 0),
			}
			for _, gw := range gws {
				if gw.DeletedAt > 0 {
					delEvent.Entities = append(delEvent.Entities, gw)
				} else {
					updateEvent.Entities = append(updateEvent.Entities, gw)
				}
			}
			if !delEvent.Empty() {
				pd.gwChan <- delEvent
			}
			if !updateEvent.Empty() {
				pd.gwChan <- updateEvent
			}
		}
		// 更新时间
		now := time.Now()
		pd.pollingAt = &now
	}
}
