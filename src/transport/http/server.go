// @Author : detaohe
// @File   : server.go
// @Description:
// @Date   : 2022/4/23 9:01 PM

package http

import (
	"gateway_kit/config"
	_ "gateway_kit/docs"
	"gateway_kit/server/endpoint"
	"gateway_kit/transport/http/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

// MakeServerHandler
// @Description:
// @return http.Handler
//
func MakeServerHandler() *gin.Engine {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		middleware.CorsMiddleware(),
		middleware.ContentTypeMiddleware(),
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/healthz"},
		}),
		middleware.AccessLogMiddleware(config.Logger),
		//middleware.RateLimitMiddleware(float64(config.All.RateLimit.Limit), config.All.RateLimit.Burst),
	)
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})

	svr := router.Group("/"+config.All.Name+"-svr", func(c *gin.Context) {
		// todo 增加一个特殊的认证
		c.Next()
	})
	pprof.RouteRegister(svr, "pprof")
	svr.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))

	svr.Use( // 超时时间
		middleware.ContextTimeout(time.Millisecond * time.Duration(config.All.Server.Timeout)),
	)

	endpoint.GatewayRouteRegister(svr)
	endpoint.HttpSvcRouteRegister(svr)
	endpoint.HttpSvcFlowRegister(svr)

	return router
}
