// @Author : detaohe
// @File   : repo
// @Description:
// @Date   : 2022/8/30 18:23

package gateway

import (
	"gateway_kit/core/lb"
	"gateway_kit/dao"
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
	// 通过redis或者数据库获取
	// todo 暂时hardcode
	svcChan    chan []*dao.HttpSvcEntity
	addrChan   chan []*lb.Node
	accessChan chan []*AccessConfig
	stopC      chan struct{}
}

func NewServiceRepo(addrC chan []*lb.Node, accessC chan []*AccessConfig) *HttpServiceRepo {
	onceRepo.Do(func() {
		RepoHttp = &HttpServiceRepo{
			svcChan:    make(chan []*dao.HttpSvcEntity),
			addrChan:   addrC,
			accessChan: accessC,
			stopC:      make(chan struct{}),
		}
		RepoHttp.Start()
	})
	return RepoHttp
}

// Start
func (repo *HttpServiceRepo) Start() {
	go func() {
	loop:
		for {
			select {
			case <-repo.stopC:
				break loop
			case entities, ok := <-repo.svcChan:
				if !ok {
					break loop
				}
				repo.updateAccess(entities)
				repo.updateAddr(entities)
			}
		}
	}()
}

func (repo *HttpServiceRepo) Stop() {
	close(repo.stopC)
	close(repo.addrChan)
	close(repo.accessChan)
}

func (repo *HttpServiceRepo) In() chan []*dao.HttpSvcEntity {
	return repo.svcChan
}

func (repo *HttpServiceRepo) updateAccess(entities []*dao.HttpSvcEntity) {
	_configs := make([]*AccessConfig, 0, len(entities))
	for _, e := range entities {
		_configs = append(_configs, &AccessConfig{
			Name:     e.Name,
			BlockIP:  e.BlockList,
			AllowIP:  e.AllowList,
			Category: ACCESS_CONTROL_SERVICE,
		})
	}
	repo.accessChan <- _configs
}

func (repo *HttpServiceRepo) updateAddr(entities []*dao.HttpSvcEntity) {
	nodes := make([]*lb.Node, 0, len(entities))
	for _, e := range entities {
		for _, addr := range e.Addrs {
			nodes = append(nodes, &lb.Node{
				SvcName: e.Name,
				Addr:    addr,
				Weight:  1,
			})
		}
	}
	repo.addrChan <- nodes
}
