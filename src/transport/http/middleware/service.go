// @Author : detaohe
// @File   : svc_name.go
// @Description:
// @Date   : 2022/9/7 18:29

package middleware

import (
	"gateway_kit/core/gateway"
	"github.com/gin-gonic/gin"
	"net/http"
)

const GwServiceName = "GwServiceName"

// 根据请求的

func ServiceNameMiddleware() gin.HandlerFunc {
	svcHandler := gateway.NewSvcMatcher()
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		host := c.Request.Host

		name, err := svcHandler.Match(host, path)
		if err != nil {
			c.JSON(http.StatusNotFound, "请求的服务不存在")
			c.Abort()
		} else {
			c.Set(GwServiceName, name)
			c.Next()
		}
	}
}
