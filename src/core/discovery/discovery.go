// @Author : detaohe
// @File   : discovery.go
// @Description:
// @Date   : 2022/9/12 09:28

package discovery

import (
	"context"
	"gateway_kit/config"
	"gateway_kit/dao"
	"gateway_kit/util"
	"go.uber.org/zap"
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
	gwChan    chan<- *dao.GwEvent
	gwInfo    *dao.GatewayEntity
	pollingAt *time.Time
	//mutex       sync.RWMutex
	//entities     []*dao.HttpSvcEntity
	//svcEntityMap map[string][]*dao.HttpSvcEntity // 以后放到redis，目前
}

func NewPollingDiscovery(svcChan chan<- *dao.SvcEvent, gwChan chan<- *dao.GwEvent) *PollingDiscovery {
	proxy := &PollingDiscovery{
		svcDao:    dao.NewHttpSvcDao(),
		gwDao:     dao.NewGatewayDao(),
		TickerSvc: util.NewTickerSvc("service-discovery", time.Minute*1, true),
		svcChan:   svcChan,
		gwChan:    gwChan,
		pollingAt: nil,
		//entities:     make([]*dao.HttpSvcEntity, 0, 1),
		//svcEntityMap: make(map[string][]*dao.HttpSvcEntity),
		//mutex:       sync.RWMutex{},
	}

	return proxy
}

func (pd *PollingDiscovery) Start() {
	pd.TickerSvc.Start(pd.endpoint())
}

func (pd *PollingDiscovery) loadServices() ([]*dao.HttpSvcEntity, error) {
	if pd.pollingAt == nil {
		return pd.svcDao.All(context.Background())
	} else {
		config.Logger.Info("loadServices", zap.Time("polling", *pd.pollingAt), zap.Time("now", time.Now()))
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
		if svcEntities, err1 := pd.loadServices(); err1 != nil {
			config.Logger.Error("load service info error", zap.Error(err1))
		} else {
			for _, entity := range svcEntities {
				config.Logger.Info("service info", zap.Any("service", entity))
				eventType := dao.EventInvalid
				switch {
				case pd.pollingAt == nil:
					eventType = dao.EventCreate
				case entity.DeletedAt > 0:
					eventType = dao.EventDelete
				case pd.pollingAt != nil && entity.CreatedAt.After(*pd.pollingAt):
					eventType = dao.EventCreate
				case pd.pollingAt != nil && entity.UpdatedAt.After(*pd.pollingAt):
					eventType = dao.EventUpdate
				default:

				}
				if eventType != dao.EventInvalid {
					pd.svcChan <- &dao.SvcEvent{
						EventType: eventType,
						Entity:    entity,
					}
				}
			}
		}

		if gws, err2 := pd.loadGateway(); err2 != nil {
			config.Logger.Error("load gateway info error", zap.Error(err2))
		} else {
			for _, gw := range gws {
				config.Logger.Info("gateway info", zap.Any("gateway", gw))
				eventType := dao.EventInvalid
				switch {
				case pd.pollingAt == nil:
					eventType = dao.EventCreate
				case pd.pollingAt != nil && gw.DeletedAt > 0:
					eventType = dao.EventDelete
				case pd.pollingAt != nil && gw.CreatedAt.After(*pd.pollingAt):
					eventType = dao.EventCreate
				case pd.pollingAt != nil && gw.UpdatedAt.After(*pd.pollingAt):
					eventType = dao.EventUpdate
				default:

				}
				if eventType != dao.EventInvalid {
					pd.gwChan <- &dao.GwEvent{
						EventType: eventType,
						Entity:    gw,
					}
				}
			}

		}
		// 更新时间
		now := time.Now().Add(-1 * time.Second * 5) // 防止竞争条件下，更新和扫描同时发生或者稍晚与扫描
		pd.pollingAt = &now
	}
}
