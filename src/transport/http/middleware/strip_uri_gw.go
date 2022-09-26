// @Author : detaohe
// @File   : strip_uri_gw.go
// @Description:
// @Date   : 2022/9/26 16:48

package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GwStripUriMiddleware(prefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, prefix := range prefixes {
			//re, _ := regexp.Compile("^/" + prefix + "/(.*)")
			//
			//urlPath := re.ReplaceAllString(c.Request.URL.Path, "$1")

			if strings.HasPrefix(c.Request.URL.Path, "/"+prefix) {
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
