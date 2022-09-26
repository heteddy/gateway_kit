// @Author : detaohe
// @File   : strip_uri_gw.go
// @Description:
// @Date   : 2022/9/26 16:48

package middleware

import (
	"gateway_kit/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

func GwStripUriMiddleware(prefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, prefix := range prefixes {
			if strings.HasPrefix(c.Request.URL.Path, "/"+prefix+"/") { //去掉指定的前缀
				pattern := "^/" + prefix + "/(.*)"
				re, _ := regexp.Compile(pattern)
				urlPath := re.ReplaceAllString(c.Request.URL.Path, "/$1")
				config.Logger.Info("strip gw uri", zap.String("urlPath", urlPath), zap.String("request.url.path", c.Request.URL.Path))
				c.Request.URL.Path = urlPath

				//// note:为了访问pprof等，如果server和gateway分开，这种设置比较有用
				//if strings.HasPrefix(c.Request.URL.Path, "/"+prefix) {
				//	c.Abort()
				//	return
				//}
				c.Next()
				break
			}
		}
		c.Next()
	}
}
