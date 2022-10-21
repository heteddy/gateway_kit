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
	// note 这里的顺序不能改，这里注册的不会使用middleware,因为注册路由的时候是拷贝了一份handlers
	admin := router.Group("/" + config.All.Name)
	admin.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))
	pprof.RouteRegister(admin, "pprof")

	router.Use( // 这里没有任何前缀
		middleware.GwStripUriMiddleware(config.All.Name),
		middleware.AccessLogMiddleware(config.Logger, "/healthz"),
		middleware.ServiceNameMiddleware(),
		middleware.IPFilterMiddleware(),
		middleware.ProtocolMiddleware(),
		middleware.ContextTimeout(time.Millisecond*time.Duration(config.All.Gateway.Timeout)),
		middleware.RateLimitMiddleware(float64(config.All.RateLimit.Limit), config.All.RateLimit.Burst),
		middleware.FlowMiddleware(),
		middleware.ReverseProxyMiddleware(lb.NewLBManager()),
	)
	return router
}
