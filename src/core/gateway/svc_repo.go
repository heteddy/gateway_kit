// @Author : detaohe
// @File   : repo
// @Description:
// @Date   : 2022/8/30 18:23

package gateway

import (
	"gateway_kit/config"
	"gateway_kit/core/lb"
	"gateway_kit/dao"
	"go.uber.org/zap"
	"sync"
)

/*
提供两种方式注册服务，
1. 适用于k8s的直接调用gateway的接口，写入client信息，gateway写入数据库并同步到redis中
2. 提供一个sdk写入到etcd, gateway通过etcd获取client的信息
*/

var RepoHttp *HttpServiceRepo
var onceRepo sync.Once

type HttpServiceRepo struct { //支持watch？
	// 接收服务发现的事件
	svcChan chan *dao.SvcEvent
	// 通知变更或者删除等
	addrChan       chan<- *lb.Node
	accessChan     chan<- *AccessConfigEvt
	svcMatcherChan chan<- *SvcMatchRule
	rateChan       chan<- *RateLimitConfigEvent
	stopC          chan struct{}
	//entities   []*dao.HttpSvcEntity
	//mutex      sync.RWMutex
}

func NewServiceRepo(
	addrC chan<- *lb.Node,
	accessC chan<- *AccessConfigEvt,
	matcherC chan<- *SvcMatchRule,
	rateC chan<- *RateLimitConfigEvent) *HttpServiceRepo {
	onceRepo.Do(func() {
		RepoHttp = &HttpServiceRepo{
			svcChan:        make(chan *dao.SvcEvent),
			addrChan:       addrC,
			accessChan:     accessC,
			svcMatcherChan: matcherC,
			rateChan:       rateC,
			stopC:          make(chan struct{}),
			//entities:   make([]*dao.HttpSvcEntity, 0, 1),
			//mutex:      sync.RWMutex{},
		}
		RepoHttp.Start()
	})
	return RepoHttp
}

// Start 启动服务
func (repo *HttpServiceRepo) Start() {
	go func() {
	loop:
		for {
			select {
			case <-repo.stopC:
				break loop
			case event, ok := <-repo.svcChan:

				if !ok {
					break loop
				}
				config.Logger.Info("receiving event", zap.Any("event", event))
				for _, entity := range event.Entities {
					repo.addrChan <- &lb.Node{
						Svc:       entity.Name,
						EventType: event.EventType,
						Addr:      entity.Addr,
						Weight:    1,
					}
					repo.accessChan <- &AccessConfigEvt{
						EventType: event.EventType,
						Name:      entity.Name,
						BlockIP:   entity.BlockList,
						AllowIP:   entity.AllowList,
						Category:  ACCESS_CONTROL_SERVICE,
					}
					repo.svcMatcherChan <- &SvcMatchRule{
						EventType: event.EventType,
						Svc:       entity.Name,
						Category:  entity.Category,
						Rule:      entity.MatchRule,
					}
					repo.rateChan <- &RateLimitConfigEvent{
						EventType: event.EventType,
						Svc:       entity.Name,
						SvcQps:    entity.ServerQps,
					}
				}
			}
		}
	}()
}

func (repo *HttpServiceRepo) Stop() {
	close(repo.stopC)
	close(repo.addrChan)
	close(repo.accessChan)
}

func (repo *HttpServiceRepo) In() chan<- *dao.SvcEvent {
	return repo.svcChan
}
