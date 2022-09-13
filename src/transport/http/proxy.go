// @Author : detaohe
// @File   : proxy.go
// @Description:
// @Date   : 2022/9/13 10:24

package http

import (
	"gateway_kit/config"
	"gateway_kit/core/lb"
	"gateway_kit/transport/http/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

func MakeProxyHandler() *gin.Engine {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		middleware.CorsMiddleware(),
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/healthz"},
		}),
	)
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})

	proxy := router.Group("/" + config.All.Name)
	//proxy.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	proxy.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))
	pprof.RouteRegister(proxy, "pprof")
	proxy.Use( // 这里没有任何前缀
		middleware.ContentTypeMiddleware(),
		middleware.AccessLogMiddleware(config.Logger),
		middleware.ContextTimeout(time.Millisecond*time.Duration(config.All.Gateway.Timeout)),
		middleware.ServiceNameMiddleware(),
		middleware.IPFilterMiddleware(),
		middleware.RateLimiter(float64(config.All.RateLimit.Limit), config.All.RateLimit.Burst),
		middleware.ReverseProxyMiddleware(lb.NewLoadBalanceMgr()),
	)

	return router
}
