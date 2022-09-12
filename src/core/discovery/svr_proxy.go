// @Author : detaohe
// @File   : svr_proxy.go
// @Description:
// @Date   : 2022/9/12 09:28

package discovery

import (
	"context"
	"gateway_kit/config"
	"gateway_kit/dao"
	"gateway_kit/util"
	"go.uber.org/zap"
	"sync"
	"time"
)

// 暂时通过admin的service 获取服务，后期可以根据数据库和

type ServiceDiscovery interface {
}

type DiscoveryProxy struct {
	svcDao *dao.HttpSvcDao // 定时获取所有的服务
	gwDao  *dao.GatewayDao
	*util.TickerSvc
	svcChan      chan<- []*dao.HttpSvcEntity
	gwChan       chan *dao.GatewayEntity
	entities     []*dao.HttpSvcEntity
	svcEntityMap map[string][]*dao.HttpSvcEntity // 以后放到redis，目前
	gwInfo       *dao.GatewayEntity
	mutex        sync.RWMutex
}

func NewDiscoveryProxy(svcChan chan []*dao.HttpSvcEntity, gwChan chan *dao.GatewayEntity) *DiscoveryProxy {
	proxy := &DiscoveryProxy{
		svcDao:       dao.NewHttpSvcDao(),
		gwDao:        dao.NewGatewayDao(),
		TickerSvc:    util.NewTickerSvc("service-discovery", time.Minute, false),
		svcChan:      svcChan,
		gwChan:       gwChan,
		entities:     make([]*dao.HttpSvcEntity, 0, 1),
		svcEntityMap: make(map[string][]*dao.HttpSvcEntity),
		mutex:        sync.RWMutex{},
	}
	proxy.Start(proxy.endpoint())
	return proxy
}

func (discovery *DiscoveryProxy) endpoint() util.SvcEndpoint {
	return func() {
		entities, err := discovery.svcDao.All(context.Background())
		if err != nil {
			config.Logger.Error("load service entities error", zap.Error(err))
			return
		}
		discovery.mutex.Lock()
		defer discovery.mutex.Unlock()
		discovery.svcEntityMap = make(map[string][]*dao.HttpSvcEntity)
		for _, e := range entities {
			if v, existed := discovery.svcEntityMap[e.Name]; existed {
				discovery.svcEntityMap[e.Name] = append(v, e)
			} else {
				_newList := make([]*dao.HttpSvcEntity, 0, 1)
				_newList = append(_newList, e)
				discovery.svcEntityMap[e.Name] = _newList
			}
		}
		discovery.entities = entities
		discovery.svcChan <- discovery.entities

		// 读取gateway
		gatewayInfos, err := discovery.gwDao.GetGateway(context.Background(), config.All.Name)
		if err != nil {
			config.Logger.Error("load gateway entities error", zap.Error(err))
			return
		}
		discovery.gwInfo = gatewayInfos[0]
		discovery.gwChan <- discovery.gwInfo
	}
}
