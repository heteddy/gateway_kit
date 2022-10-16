// @Author : detaohe
// @File   : ratelimit.go
// @Description:
// @Date   : 2022/9/12 21:09

package gateway

import (
	"gateway_kit/config"
	"gateway_kit/dao"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"sync"
)

var onceLimiter sync.Once
var Limiter *RateLimiter

type RateLimitConfigEvent struct {
	EventType int
	Svc       string
	SvcQps    int
}
type RateLimiter struct {
	mutex      sync.RWMutex
	stopC      chan struct{}
	configC    chan *RateLimitConfigEvent
	svcQps     map[string]*RateLimitConfigEvent
	svcLimiter map[string]*rate.Limiter // 针对每个服务创建一个limiter，这里limiter需要改成基于redis
}

func NewRateLimiter() *RateLimiter {
	onceLimiter.Do(func() {
		Limiter = &RateLimiter{
			mutex:      sync.RWMutex{},
			stopC:      make(chan struct{}),
			configC:    make(chan *RateLimitConfigEvent),
			svcQps:     make(map[string]*RateLimitConfigEvent),
			svcLimiter: make(map[string]*rate.Limiter),
		}
	})
	return Limiter
}

func (rl *RateLimiter) update(c *RateLimitConfigEvent) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	config.Logger.Info("update RateLimiter", zap.Any("RateLimitConfigEvent", c))
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
				config.Logger.Warn("RateLimiter exit")
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
func (rl *RateLimiter) In() chan<- *RateLimitConfigEvent {
	return rl.configC
}

func (rl *RateLimiter) Allow(name string) bool {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()
	if limiter, existed := rl.svcLimiter[name]; existed {
		return limiter.Allow()
	} else {
		return true // 默认不限流
	}
}
