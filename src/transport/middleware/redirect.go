// @Author : detaohe
// @File   : redirect.go
// @Description:
// @Date   : 2022/9/6 20:57

package middleware

import "github.com/gin-gonic/gin"

func RedirectMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
