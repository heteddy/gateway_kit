// @Author : detaohe
// @File   : repo
// @Description:
// @Date   : 2022/8/30 18:23

package gateway

import (
	"gateway_kit/gateway/lb"
	"sync"
)

/*
提供两种方式注册服务，
1. 适用于k8s的直接调用gateway的接口，写入client信息，gateway写入数据库并同步到redis中
2. 提供一个sdk写入到etcd, gateway通过etcd获取client的信息
*/

type ServiceRepo struct { //支持watch？
	// 通过redis或者数据库获取
	// todo 暂时hardcode
	serviceAddrs map[string][]string
	listenerMap  map[string][]lb.Listener
	mutex        sync.Mutex //
	// 启动一个goroutine，监控数据变化

}

func NewServiceRepo() *ServiceRepo {
	serviceAddrs := make(map[string][]string)
	serviceAddrs["proxy"] = []string{
		"192.168.64.7:9192",
		"192.168.64.7:9193",
	}
	return &ServiceRepo{
		serviceAddrs: serviceAddrs,
		listenerMap:  make(map[string][]lb.Listener),
	}
}

// Watch 监控服务的变化
func (repo *ServiceRepo) Watch(svcName string, l lb.Listener) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	if ls, existed := repo.listenerMap[svcName]; existed {
		ls = append(ls, l)
		repo.listenerMap[svcName] = ls
	} else {
		ls := make([]lb.Listener, 0, 1)
		ls = append(ls, l)
		repo.listenerMap[svcName] = ls
	}
}

// Update 更新地址
func (repo *ServiceRepo) Update(name string, addrs []string) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.serviceAddrs[name] = addrs
	if ls, existed := repo.listenerMap[name]; existed {
		for _, l := range ls {
			l.Update(name, addrs)
		}
	}
}
