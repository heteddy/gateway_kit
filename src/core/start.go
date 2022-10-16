// @Author : detaohe
// @File   : start.go
// @Description:
// @Date   : 2022/9/13 12:07

package core

import (
	"fmt"
	"gateway_kit/config"
	"gateway_kit/core/discovery"
	"gateway_kit/core/gateway"
	"gateway_kit/core/gateway/flow"
	"gateway_kit/core/lb"
)

func Start() {
	ac := gateway.NewAccessController() //接收黑白名单
	ac.Start()
	bm := lb.NewLBManager() // 接收服务地址
	bm.Start()
	matcher := gateway.NewSvcMatcher()
	matcher.Start()
	rate := gateway.NewRateLimiter()
	rate.Start()
	proto := gateway.NewProtocolTransCtrl()
	proto.Start()
	repo := gateway.NewServiceRepo(bm.In(), ac.In(), matcher.In(), rate.In(), proto.In()) // 接收服务配置信息
	repo.Start()
	gw := gateway.NewGwController(ac.In()) // 接收gateway配置信息
	gw.Start()
	pd := discovery.NewPollingDiscovery(repo.In(), gw.In())
	pd.Start()

	flowC := flow.NewFlowCollector(config.All.Name)
	flowC.Start()

	// 服务发现，目前轮训数据库
	fmt.Printf("gateway started......\n")
}
