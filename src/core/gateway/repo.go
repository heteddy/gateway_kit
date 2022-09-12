// @Author : detaohe
// @File   : repo
// @Description:
// @Date   : 2022/8/30 18:23

package gateway

import (
	"gateway_kit/lb"
	"sync"
)

/*
提供两种方式注册服务，
1. 适用于k8s的直接调用gateway的接口，写入client信息，gateway写入数据库并同步到redis中
2. 提供一个sdk写入到etcd, gateway通过etcd获取client的信息
*/

type HttpServiceRepo struct { //支持watch？
	// 通过redis或者数据库获取
	// todo 暂时hardcode
	serviceAddrs    map[string][]string
	addrListenerMap map[string][]lb.Listener
	mutex           sync.Mutex //
	// 启动一个goroutine，监控数据变化
	blockListListener map[string][]
}

func NewServiceAddrRepo() *HttpServiceRepo {
	serviceAddrs := make(map[string][]string)
	serviceAddrs["proxy"] = []string{
		"192.168.64.7:9192",
		"192.168.64.7:9193",
	}
	return &HttpServiceRepo{
		serviceAddrs:    serviceAddrs,
		addrListenerMap: make(map[string][]lb.Listener),
	}
}

// Watch 监控服务的变化
func (repo *HttpServiceRepo) Watch(svcName string, l lb.Listener) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	if ls, existed := repo.addrListenerMap[svcName]; existed {
		ls = append(ls, l)
		repo.addrListenerMap[svcName] = ls
	} else {
		ls := make([]lb.Listener, 0, 1)
		ls = append(ls, l)
		repo.addrListenerMap[svcName] = ls
	}
}

// UpdateAddr 更新地址，因为地址和block list可能需要分开处理
func (repo *HttpServiceRepo) UpdateAddr(name string, addrs []string) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.serviceAddrs[name] = addrs
	if ls, existed := repo.addrListenerMap[name]; existed {
		for _, l := range ls {
			l.Update(name, addrs)
		}
	}
}
