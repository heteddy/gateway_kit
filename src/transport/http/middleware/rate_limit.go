// @Author : detaohe
// @File   : rate_limit.go
// @Description:
// @Date   : 2022/9/5 16:34

package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

func RateLimiter(limit float64, burst int) gin.HandlerFunc {
	// note 这里只能是单机版
	limiter := rate.NewLimiter(rate.Limit(limit), burst)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
		} else {
			// todo 这里需要基于redis，否则多实例的时候
			// 是不是可以直接打印
			//c.Writer.Write([]byte(fmt.Sprintf("ratelimit limit=%v,burst=%v\n", limiter.Limit(), limiter.Burst())))
			c.JSON(http.StatusTooManyRequests, fmt.Sprintf("ratelimit limit=%v,burst=%v\n", limiter.Limit(), limiter.Burst()))
			c.Abort()
		}
	}
}
