// @Author : detaohe
// @File   : loadbalance.go
// @Description:
// @Date   : 2022/4/23 8:32 PM

package lb

import (
	"math/rand"
	"time"
)

type LoadBalancer interface {
	// GetService 通过一个svcName 获取真实的地址
	GetService([]string) (string, error)
}

type weightRoundRobinLB struct {
}

func (lb *weightRoundRobinLB) GetService(hosts []string) (string, error) {
	return hosts[0], nil
}

type randLB struct {
}

func NewRandomLB() LoadBalancer {
	return &randLB{}
}

func (lb *randLB) GetService(hosts []string) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Int63()
	idx := n % int64(len(hosts))
	return hosts[idx], nil
}
