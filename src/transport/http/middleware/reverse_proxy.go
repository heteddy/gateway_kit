// @Author : detaohe
// @File   : reverse_proxy.go
// @Description:
// @Date   : 2022/9/13 10:31

package middleware

import (
	"gateway_kit/config"
	"gateway_kit/core/gateway"
	"gateway_kit/core/lb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReverseProxyMiddleware(balanceMgr *lb.LoadBalanceMgr) gin.HandlerFunc {
	return func(c *gin.Context) {
		if svc, existed := c.Get(KeyGwSvcName); existed {
			svcName := svc.(string)
			builder := gateway.NewProxyBuilder()
			//balancer := lb.NewRoundRobin() // 根据配置
			var scheme string
			if _scheme, exist2 := c.Get(KeySvcRequestScheme); exist2 {
				scheme = _scheme.(string)
			} else {
				config.Logger.Warn("scheme not found(ReverseProxyMiddleware)")
			}
			// 根据配置获取
			_proxy := builder.BuildHttpProxy(balanceMgr.Get(lb.LbRandom), svcName, scheme)
			_proxy.ServeHTTP(c.Writer, c.Request)
			c.Abort() // proxy的server是一个假的路由，因此这里必须加上Abort
		} else {
			c.JSON(http.StatusNotFound, "请求的服务不存在(ReverseProxyMiddleware)")
			c.Abort()
		}
	}
}
