// @Author : detaohe
// @File   : svc_name.go
// @Description:
// @Date   : 2022/9/7 18:29

package middleware

import (
	"gateway_kit/config"
	"gateway_kit/core/gateway"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const KeyGwSvcName = "KeyGwSvcName"

// 根据请求的

func ServiceNameMiddleware() gin.HandlerFunc {
	svcHandler := gateway.NewSvcMatcher()
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		host := c.Request.Host
		config.Logger.Info("strip gw uri", zap.String("host", host), zap.String("request.url.path", c.Request.URL.Path))
		name, err := svcHandler.Match(host, path)
		if err != nil {
			c.JSON(http.StatusNotFound, "请求的服务不存在(ServiceNameMiddleware)")
			c.Abort()
		} else {
			c.Set(KeyGwSvcName, name)
			c.Next()
		}
	}
}
