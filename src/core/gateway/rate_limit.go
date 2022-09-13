// @Author : detaohe
// @File   : rate_limit.go
// @Description:
// @Date   : 2022/9/12 21:09

package gateway

import (
	"golang.org/x/time/rate"
	"sync"
)

var onceLimiter sync.Once
var Limiter *RateLimiter

type RateLimiter struct {
	mutex   sync.RWMutex
	stopC   chan struct{}
	svrQps  map[string]int64
	limiter map[string]*rate.Limiter
}

func NewRateLimiter() *RateLimiter {
	onceLimiter.Do(func() {
		Limiter = &RateLimiter{
			mutex:   sync.RWMutex{},
			stopC:   make(chan struct{}),
			svrQps:  make(map[string]int64),
			limiter: make(map[string]*rate.Limiter),
		}
	})
	return Limiter
}
