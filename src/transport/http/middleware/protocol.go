// @Author : detaohe
// @File   : protocol.go
// @Description:
// @Date   : 2022/9/25 11:44

package middleware

import (
	"gateway_kit/core/gateway"
	"github.com/gin-gonic/gin"
)

const (
	KeySvcRequestScheme = "KeySvcRequestScheme"
	//KeySvcRequestHost   = "KeySvcRequestHost"
)

func ProtocolMiddleware() gin.HandlerFunc {
	protocolCtrl := gateway.NewProtocolTransCtrl()
	return func(c *gin.Context) {
		if svc, existed := c.Get(KeyGwSvcName); !existed {
			// 这里不应该不存在，因为service中间件应该会拒绝掉
			svcName := svc.(string)
			isHttps := protocolCtrl.IsHttps(svcName)
			isWebSocket := protocolCtrl.IsWebsocket(svcName)

			switch {
			case isHttps && isWebSocket:
				c.Set(KeySvcRequestScheme, "wss")
			case isHttps:
				c.Set(KeySvcRequestScheme, "https")
			case isWebSocket:
				c.Set(KeySvcRequestScheme, "ws")
			default:
				c.Set(KeySvcRequestScheme, "http")
			}
		}
		c.Next()
	}
}
