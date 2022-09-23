// @Author : detaohe
// @File   : rate_limit.go
// @Description:
// @Date   : 2022/9/12 21:09

package gateway

import (
	"gateway_kit/dao"
	"golang.org/x/time/rate"
	"sync"
)

var onceLimiter sync.Once
var Limiter *RateLimiter

type RateLimitConfig struct {
	EventType int
	Svc       string
	SvcQps    int
}
type RateLimiter struct {
	mutex      sync.RWMutex
	stopC      chan struct{}
	configC    chan *RateLimitConfig
	svcQps     map[string]*RateLimitConfig
	svcLimiter map[string]*rate.Limiter // 针对每个服务创建一个limiter，这里limiter需要改成基于redis
}

func NewRateLimiter() *RateLimiter {
	onceLimiter.Do(func() {
		Limiter = &RateLimiter{
			mutex:      sync.RWMutex{},
			stopC:      make(chan struct{}),
			configC:    make(chan *RateLimitConfig),
			svcQps:     make(map[string]*RateLimitConfig),
			svcLimiter: make(map[string]*rate.Limiter),
		}
	})
	return Limiter
}

func (rl *RateLimiter) update(c *RateLimitConfig) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	if c.EventType == dao.EventDelete {
		delete(rl.svcQps, c.Svc)
		delete(rl.svcLimiter, c.Svc)
	} else {
		rl.svcQps[c.Svc] = c
		rl.svcLimiter[c.Svc] = rate.NewLimiter(rate.Limit(c.SvcQps), c.SvcQps*2) // 每秒的个数和桶大小
	}
}

func (rl *RateLimiter) run() {
loop:
	for {
		select {
		case <-rl.stopC:
			break loop
		case c, ok := <-rl.configC:
			if !ok {
				break loop
			}
			rl.update(c)
		}
	}
}

func (rl *RateLimiter) Start() {
	go rl.run()
}

func (rl *RateLimiter) Stop() {
	close(rl.stopC)
}
func (rl *RateLimiter) In() chan<- *RateLimitConfig {
	return rl.configC
}

func (rl *RateLimiter) Allow(name string) bool {
	if limiter, existed := rl.svcLimiter[name]; existed {
		return limiter.Allow()
	} else {
		return true // 默认不限流
	}
}
