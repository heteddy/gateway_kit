// @Author : detaohe
// @File   : svc_name.go
// @Description:
// @Date   : 2022/9/12 20:08

package gateway

import (
	"errors"
	"gateway_kit/config"
	"gateway_kit/dao"
	"go.uber.org/zap"
	"strings"
	"sync"
)

var onceMatcher sync.Once
var svcMatcher *SvcMatcher

type SvcMatchRule struct {
	Svc       string
	EventType int
	Category  int
	Rule      string
}

func (r SvcMatchRule) Is(other *SvcMatchRule) bool {
	return r.Svc == other.Svc && r.Category == other.Category && r.Rule == other.Rule
}

func (r SvcMatchRule) Match(host, path string) bool {
	switch r.Category {
	case dao.SvcCategoryUrlPrefix:
		// todo 增加正则表达式
		if strings.HasPrefix(path, r.Rule) {
			return true
		}
	case dao.SvcCategoryHost:
		if host == r.Rule {
			return true
		}
	default:
		//
		config.Logger.Warn(
			"error of rule category",
		)

	}
	return false
}

type SvcMatcher struct {
	mutex    sync.RWMutex
	svcRules []*SvcMatchRule
	stopC    chan struct{}
	ruleC    chan *SvcMatchRule
}

func NewSvcMatcher() *SvcMatcher {
	onceMatcher.Do(func() {
		svcMatcher = &SvcMatcher{
			mutex:    sync.RWMutex{},
			svcRules: make([]*SvcMatchRule, 0, 10),
			stopC:    make(chan struct{}),
			ruleC:    make(chan *SvcMatchRule),
		}
		//svcMatcher.Start()

	})
	return svcMatcher
}

func (matcher *SvcMatcher) Match(host, path string) (string, error) {
	matcher.mutex.RLock()
	defer matcher.mutex.RUnlock()
	for _, entity := range matcher.svcRules {
		if entity.Match(host, path) {
			return entity.Svc, nil
		}
	}
	return "", errors.New("service not found")
}

func (matcher *SvcMatcher) In() chan<- *SvcMatchRule {
	return matcher.ruleC
}

func (matcher *SvcMatcher) delete(rule *SvcMatchRule) {
	for idx, r := range matcher.svcRules {
		if r.Is(rule) {
			matcher.svcRules = append(matcher.svcRules[0:idx], matcher.svcRules[idx+1:]...)
		}
	}
}
func (matcher *SvcMatcher) update(rule *SvcMatchRule) {
	matcher.mutex.Lock()
	defer matcher.mutex.Unlock()
	config.Logger.Info("update svc matcher", zap.Any("rule", rule))
	switch rule.EventType {
	case dao.EventDelete:
		matcher.delete(rule)
	case dao.EventUpdate:
		matcher.delete(rule)
		matcher.svcRules = append(matcher.svcRules, rule)
	case dao.EventCreate:
		matcher.svcRules = append(matcher.svcRules, rule)
	default:

	}
}

func (matcher *SvcMatcher) runLoop() {
loop:
	for {
		select {
		case rule, ok := <-matcher.ruleC:
			if !ok {
				break loop
			}
			matcher.update(rule)
		case <-matcher.stopC:
			break loop
		}
	}
}

func (matcher *SvcMatcher) Start() {
	go matcher.runLoop()
}

func (matcher *SvcMatcher) Stop() {
	close(matcher.stopC)
}
