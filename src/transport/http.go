// @Author : detaohe
// @File   : http.go
// @Description:
// @Date   : 2022/4/23 9:01 PM

package transport

import (
	"gateway_kit/config"
	"gateway_kit/endpoint"
	"gateway_kit/transport/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// MakeHttpHandler
// @Description:
// @return http.Handler
//
func MakeHttpHandler() *gin.Engine {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		middleware.CorsMiddleware(),
		middleware.ContentTypeMiddleware(),
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/healthz"},
		}),
		middleware.AccessLogMiddleware(config.Logger),
		middleware.RateLimiter(float64(config.All.RateLimit.Limit), config.All.RateLimit.Burst),
	)

	servicePath := router.Group("/"+config.All.Service, func(c *gin.Context) {
		// todo 增加一个特殊的认证
		c.Next()
	})
	pprof.RouteRegister(servicePath, "pprof")
	endpoint.StringRouteReg(servicePath)
	return router
}
