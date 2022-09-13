// @Author : detaohe
// @File   : reverse_proxy.go
// @Description:
// @Date   : 2022/9/13 10:31

package middleware

import (
	"gateway_kit/core/gateway"
	"gateway_kit/core/lb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReverseProxyMiddleware(balanceMgr *lb.LoadBalanceMgr) gin.HandlerFunc {
	return func(c *gin.Context) {
		if svc, existed := c.Get(GwServiceName); existed {
			svcName := svc.(string)
			builder := gateway.NewProxyBuilder()
			//balancer := lb.NewRoundRobin() // 根据配置
			_proxy := builder.BuildHttpProxy(balanceMgr.Get(lb.Lb_Random), svcName)
			_proxy.ServeHTTP(c.Writer, c.Request)
			c.Abort() // proxy的server是一个假的路由，因此这里必须加上Abort
			return
		} else {
			c.JSON(http.StatusNotFound, "请求的服务不存在")
			c.Abort()
			return
		}
	}
}
