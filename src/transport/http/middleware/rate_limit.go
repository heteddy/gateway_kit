// @Author : detaohe
// @File   : rate_limit.go
// @Description:
// @Date   : 2022/9/5 16:34

package middleware

import (
	"fmt"
	"gateway_kit/core/gateway"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

func RateLimiteMiddleware(limit float64, burst int) gin.HandlerFunc {
	// note 这里只能是单机版
	sysLimiter := rate.NewLimiter(rate.Limit(limit), burst)
	svcLimiter := gateway.NewRateLimiter()
	return func(c *gin.Context) {
		if svc, existed := c.Get(KeyGwSvcName); !existed {
			// 这里不应该不存在，因为service中间件应该会拒绝掉
			c.Abort()
		} else {
			svcName := svc.(string)
			if !sysLimiter.Allow() {
				c.JSON(http.StatusTooManyRequests, fmt.Sprintf("gateway reject your request ratelimit limit=%v,burst=%v\n", sysLimiter.Limit(), sysLimiter.Burst()))
				c.Abort()
				return
			}
			if !svcLimiter.Allow(svcName) {
				c.JSON(http.StatusTooManyRequests, fmt.Sprintf("%s reject your request\n", svc))
				c.Abort()
				return
			}
			c.Next()
		}
	}
}
