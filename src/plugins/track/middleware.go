// @Author : detaohe
// @File   : middleware.go
// @Description:
// @Date   : 2022/10/21 19:13

package track

import (
	"gateway_kit/config"
	"github.com/gin-gonic/gin"
)

func EventTrackMiddleware(c config.KafkaSinkConfig) gin.HandlerFunc {
	tracker := NewEventTracker(c)
	return func(c *gin.Context) {
		tracker.Track(c.Request)
		c.Next()
	}
}
