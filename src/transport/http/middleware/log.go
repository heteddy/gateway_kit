// @Author : detaohe
// @File   : access
// @Description:
// @Date   : 2022/9/4 17:56

package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type AccessLogger struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogger) Write(body []byte) (int, error) {
	if n, err := w.body.Write(body); err != nil {
		return n, err
	} else {
		return w.ResponseWriter.Write(body)
	}
}

func AccessLogMiddleware(logger *zap.Logger, ignores ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyWriter := &AccessLogger{
			body:           bytes.NewBufferString(""),
			ResponseWriter: ctx.Writer,
		}
		begin := time.Now()
		ctx.Next()
		var ignore bool
		for _, p := range ignores {
			if strings.HasPrefix(ctx.Request.URL.Path, p) {
				ignore = true
				break
			}
		}
		if !ignore {
			method := logger.Info
			if bodyWriter.ResponseWriter.Status() >= http.StatusBadRequest {
				method = logger.Warn
			}
			end := time.Now()
			delta := end.Sub(begin)
			method("access log",
				zap.String("method", ctx.Request.Method),
				zap.String("path", ctx.Request.URL.Path),
				zap.Int("status", bodyWriter.ResponseWriter.Status()),
				zap.Time("begin", begin),
				zap.Time("end", end),
				zap.Duration("duration", delta),
			)
		}
	}
}
