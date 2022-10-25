// @Author : detaohe
// @File   : start.go
// @Description:
// @Date   : 2022/9/13 12:07

package core

import (
	"fmt"
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
	rewriter := gateway.NewRewriter()
	rewriter.Start()
	repo := gateway.NewServiceRepo(bm.In(), ac.In(), matcher.In(), rate.In(), proto.In(), rewriter.In()) // 接收服务配置信息
	repo.Start()
	gw := gateway.NewGwController(ac.In()) // 接收gateway配置信息
	gw.Start()
	pd := discovery.NewPollingDiscovery(repo.In(), gw.In())
	pd.Start()

	flowC := flow.NewFlowCollector()
	flowC.Start()

	// note 如果需要配置埋点在这里加上
	//tracker := track.NewEventTracker(config.All.KafkaEventProducer)
	//tracker.Start()

	// 服务发现，目前轮训数据库
	fmt.Printf("gateway started......\n")
}
