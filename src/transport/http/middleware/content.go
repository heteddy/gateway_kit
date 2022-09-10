// @Author : detaohe
// @File   : content
// @Description:
// @Date   : 2022/9/4 17:48

package middleware

import (
	"github.com/gin-gonic/gin"
)

func ContentTypeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json;charset=utf-8")
		ctx.Next()
	}
}
