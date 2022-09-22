// @Author : detaohe
// @File   : svc_name.go
// @Description:
// @Date   : 2022/9/12 20:08

package gateway

import (
	"errors"
	"gateway_kit/config"
	"gateway_kit/dao"
	"strings"
	"sync"
)

var onceMatcher sync.Once
var Matcher *SvcMatcher

type SvcMatchRule struct {
	Svc       string
	EventType int
	Category  int
	Rule      string
}

type SvcMatcher struct {
	mutex    sync.RWMutex
	svcRules []*SvcMatchRule
	stopC    chan struct{}
	ruleC    chan []*SvcMatchRule
}

func NewSvcMatcher() *SvcMatcher {
	onceMatcher.Do(func() {
		Matcher = &SvcMatcher{
			mutex:    sync.RWMutex{},
			svcRules: nil,
			stopC:    make(chan struct{}),
			ruleC:    make(chan []*SvcMatchRule),
		}
		Matcher.Start()

	})
	return Matcher
}

func (svc *SvcMatcher) Match(host, path string) (string, error) {
	svc.mutex.RLock()
	defer svc.mutex.RUnlock()
	for _, entity := range svc.svcRules {
		switch entity.Category {
		case dao.SvcCategoryUrlPrefix:
			// todo 增加正则表达式
			if strings.HasPrefix(path, entity.Rule) {
				return entity.Svc, nil
			}
		case dao.SvcCategoryHost:
			if host == entity.Rule {
				return entity.Svc, nil
			}
		default:
			//
			config.Logger.Warn(
				"error of entity category",
			)

		}
	}
	return "", errors.New("service not found")
}

func (svc *SvcMatcher) In() chan []*SvcMatchRule {
	return svc.ruleC
}
func (svc *SvcMatcher) runLoop() {
loop:
	for {
		select {
		case rules, ok := <-svc.ruleC:
			if !ok {
				break loop
			}
			svc.svcRules = rules
		case <-svc.stopC:
			break loop
		}
	}
}

func (svc *SvcMatcher) Start() {
	go svc.runLoop()
}

func (svc *SvcMatcher) Stop() {
	close(svc.stopC)
}
