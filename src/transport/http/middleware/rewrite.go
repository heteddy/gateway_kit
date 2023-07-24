// @Author : detaohe
// @File   : rewrite.go
// @Description:
// @Date   : 2022/9/8 21:07

package middleware

import (
	"gateway_kit/config"
	"gateway_kit/core/gateway"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

/*
rewrite和redirect的区别
redirect是在客户端的角度，客户端发送到服务器，服务器返回301和重定向的地址，客户端自动请求新地址
rewrite 是在服务器的角度 /resource 重写到 /different-resource 时，客户端会请求 /resource ，
并且服务器会在内部提取 /different-resource 处的资源。尽管客户端可能能够检索已重写URL处的资源
*/

func RewriteUriMiddleware() gin.HandlerFunc {
	rewriter := gateway.NewRewriter()
	return func(c *gin.Context) {
		if svc, existed := c.Get(KeyGwSvcName); existed {
			svcName := svc.(string)
			newPath, err := rewriter.Rewrite(svcName, c.Request.URL.Path)
			if err != nil {
				config.Logger.Error("rewrite failure:", zap.Error(err))
			} else {
				c.Request.URL.Path = newPath
			}
			c.Next()
		} else {
			c.JSON(http.StatusNotFound, "no service found")
			c.Abort()
		}

	}
}
