// @Author : detaohe
// @File   : timeout.go
// @Description:
// @Date   : 2022/9/8 16:26

package middleware

import (
	"context"
	"errors"
	"gateway_kit/core/gateway"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func ContextTimeout(t time.Duration) gin.HandlerFunc {
	protocolCtrl := gateway.NewProtocolTransCtrl()
	return func(c *gin.Context) {
		if svc, existed := c.Get(KeyGwSvcName); existed {
			// 这里不应该不存在，因为service中间件应该会拒绝掉
			svcName := svc.(string)
			isWebSocket := protocolCtrl.IsWebsocket(svcName)
			if isWebSocket {

			} else {
				ctx, cancel := context.WithTimeout(c.Request.Context(), t)
				defer func() {
					if errors.Is(ctx.Err(), context.DeadlineExceeded) {
						c.Writer.WriteHeader(http.StatusGatewayTimeout)
						c.Abort()
					}
					cancel()
				}()
				c.Request = c.Request.WithContext(ctx)
			}
		}
		c.Next()
	}
}
