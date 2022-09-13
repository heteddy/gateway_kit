// @Author : detaohe
// @File   : start.go
// @Description:
// @Date   : 2022/9/13 12:07

package core

import (
	"gateway_kit/core/discovery"
	"gateway_kit/core/gateway"
	"gateway_kit/core/lb"
)

func InitCore() {
	ac := gateway.NewAccessController()              //接收黑白名单
	bm := lb.NewLoadBalanceMgr()                     // 接收服务地址
	repo := gateway.NewServiceRepo(bm.In(), ac.In()) // 接收服务配置信息
	gw := gateway.NewGwConfig(ac.In())               // 接收gateway配置信息
	discovery.NewDiscoveryProxy(repo.In(), gw.In())  // 服务发现，目前轮训数据库
}
