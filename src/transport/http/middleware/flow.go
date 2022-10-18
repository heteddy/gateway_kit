// @Author : detaohe
// @File   : flow.go
// @Description:
// @Date   : 2022/10/17 21:12

package middleware

import (
	"gateway_kit/core/gateway/flow"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FlowMiddleware() gin.HandlerFunc {
	flowHandler := flow.NewFlowCollector()
	return func(c *gin.Context) {
		if svc, existed := c.Get(KeyGwSvcName); existed {
			svcName := svc.(string)
			flowHandler.Call(svcName, c.Request.URL.Path, c.Request.Method)
			c.Next()
		} else {
			c.JSON(http.StatusNotFound, "请求的服务不存在(IPFilterMiddleware)")
			c.Abort()
		}

	}
}
