// @Author : detaohe
// @File   : http.go
// @Description:
// @Date   : 2022/4/23 9:01 PM

package http

import (
	"gateway_kit/config"
	"gateway_kit/endpoint"
	"gateway_kit/transport/http/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"time"
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

	admin := router.Group("/"+config.All.Name+"/admin", func(c *gin.Context) {
		// todo 增加一个特殊的认证
		c.Next()
	})
	admin.Use( // 超时时间
		middleware.ContextTimeout(time.Millisecond * time.Duration(config.All.Server.Timeout)),
	)
	pprof.RouteRegister(admin, "pprof")
	endpoint.StringRouteReg(admin)

	proxy := router.Group("/" + config.All.Name + "/proxy")
	proxy.Use(
		middleware.ContextTimeout(time.Millisecond*time.Duration(config.All.Gateway.Timeout)),
		middleware.ServiceNameMiddleware(),
	)
	return router
}
