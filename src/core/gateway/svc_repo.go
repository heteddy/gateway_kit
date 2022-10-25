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
	protoChan      chan<- *ProtocolSupportedEvt
	rewriteChan    chan<- *SvcRewriteRule
	stopC          chan struct{}

	//entities   []*dao.HttpSvcEntity
	//mutex      sync.RWMutex
}

func NewServiceRepo(
	addrC chan<- *lb.Node,
	accessC chan<- *AccessConfigEvt,
	matcherC chan<- *SvcMatchRule,
	rateC chan<- *RateLimitConfigEvent,
	protoC chan<- *ProtocolSupportedEvt,
	rewriteC chan<- *SvcRewriteRule) *HttpServiceRepo {
	onceRepo.Do(func() {
		RepoHttp = &HttpServiceRepo{
			svcChan:        make(chan *dao.SvcEvent),
			addrChan:       addrC,
			accessChan:     accessC,
			svcMatcherChan: matcherC,
			rateChan:       rateC,
			protoChan:      protoC,
			rewriteChan:    rewriteC,
			stopC:          make(chan struct{}),
		}
	})
	return RepoHttp
}

func (repo *HttpServiceRepo) runLoop() {
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
			entity := event.Entity
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

			repo.protoChan <- &ProtocolSupportedEvt{
				EventType:   event.EventType,
				Svc:         entity.Name,
				IsWebsocket: entity.IsWebsocket,
				IsHttps:     entity.IsHttps,
			}
			repo.rewriteChan <- &SvcRewriteRule{
				EventType:   event.EventType,
				Svc:         entity.Name,
				RewriteUrls: entity.UrlRewrite,
				Patterns:    make([]rewritePattern, 0, len(entity.UrlRewrite)),
			}

		}
	}
}

// Start 启动服务
func (repo *HttpServiceRepo) Start() {
	go repo.runLoop()
}

func (repo *HttpServiceRepo) Stop() {
	close(repo.stopC)
	close(repo.addrChan)
	close(repo.accessChan)
	close(repo.svcMatcherChan)
	close(repo.rateChan)
	close(repo.protoChan)
}

func (repo *HttpServiceRepo) In() chan<- *dao.SvcEvent {
	return repo.svcChan
}
