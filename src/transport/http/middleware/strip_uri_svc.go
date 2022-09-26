// @Author : detaohe
// @File   : strip_uri_svc.go
// @Description:
// @Date   : 2022/9/8 21:21

package middleware

import (
	"gateway_kit/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

//mockbin	true	/mockbin/some_path	/some_path
func StripUriMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if svc, existed := c.Get(KeyGwSvcName); existed {
			svcName := svc.(string)
			config.Logger.Info("strip url", zap.String("servicename", svcName))
			c.Next()
		} else {
			c.JSON(http.StatusNotFound, "请求的服务不存在(StripUriMiddleware)")
			c.Abort()
		}
	}
}
