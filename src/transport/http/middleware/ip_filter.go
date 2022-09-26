// @Author : detaohe
// @File   : ip_filter.go
// @Description:
// @Date   : 2022/9/8 20:45

package middleware

import (
	"gateway_kit/core/gateway"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IPFilterMiddleware() gin.HandlerFunc {
	filterHandler := gateway.NewAccessController()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if svc, existed := c.Get(KeyGwSvcName); existed {
			svcName := svc.(string)
			if filterHandler.Allow(svcName, ip) {
				c.Next()
			} else {
				c.JSON(http.StatusMethodNotAllowed, "无权访问(IPFilterMiddleware)")
				c.Abort()
			}
		} else {
			c.JSON(http.StatusNotFound, "请求的服务不存在(IPFilterMiddleware)")
			c.Abort()
		}

	}
}
