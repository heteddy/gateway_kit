// @Author : detaohe
// @File   : cors.
// @Description:
// @Date   : 2022/9/4 16:35

package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//[]string{http.MethodPost, http.MethodGet, http.MethodOptions, http.MethodPut, http.MethodDelete, http.MethodConnect, http.MethodPatch, http.MethodHead, http.MethodTrace},
		method := ctx.Request.Method
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", strings.Join([]string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"AccessToken",
			"X-CSRF-Token",
			"Authorization",
			"Token",
			"x-token"}, ","))
		ctx.Header("Access-Control-Expose-Headers", strings.Join([]string{
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Content-Length",
			"Content-Type",
			"Last-Modified",
			"Expires",
		}, ","))
		ctx.Header("Access-Control-Allow-Credentials", "false")
		ctx.Header("Access-Control-Max-Age", "86400") // 24*3600
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, gin.H{})
			return
		}
		ctx.Next()
	}
}
