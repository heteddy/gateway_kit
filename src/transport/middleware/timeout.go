// @Author : detaohe
// @File   : timeout.go
// @Description:
// @Date   : 2022/9/8 16:26

package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func ContextTimeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer func() {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}
			cancel()
		}()
		
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
