// @Author : detaohe
// @File   : rewrite.go
// @Description:
// @Date   : 2022/10/25 16:48

package gateway

import (
	"errors"
	"gateway_kit/config"
	"gateway_kit/dao"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"sync"
)

var onceRewriter sync.Once
var svcRewriter *SvcRewriter

type rewritePattern struct {
	Origin      string
	Destination string
}

type SvcRewriteRule struct {
	Svc         string
	EventType   int
	RewriteUrls []string
	Patterns    []rewritePattern
}

func (r *SvcRewriteRule) genPattern() error {
	for _, rewrites := range r.RewriteUrls {
		ss := strings.Split(rewrites, " ")
		if len(ss) < 3 {
			return errors.New("pattern error\nexample: rewrite /api/v1/dict/(.*)  /api/v1/$1")
		}
		if ss[0] != "rewrite" {
			return errors.New("pattern error\nexample: rewrite /api/v1/dict/(.*)  /api/v1/$1")
		}
		r.Patterns = append(r.Patterns, rewritePattern{
			Origin:      ss[1],
			Destination: ss[2],
		})
	}

	return nil
}

func (r SvcRewriteRule) Is(other *SvcRewriteRule) bool {
	return r.Svc == other.Svc
}

func (r SvcRewriteRule) Rewrite(path string) (string, error) {
	for _, op := range r.Patterns {
		re, err := regexp.Compile(op.Origin)
		if err != nil {
			return "", err
		}
		if re.MatchString(path) {
			dst := re.ReplaceAllString(path, op.Destination)
			return dst, nil
		}
	}
	return "", errors.New("no match rewrite rule")
}

type SvcRewriter struct {
	mutex    sync.RWMutex
	svcRules []*SvcRewriteRule
	stopC    chan struct{}
	ruleC    chan *SvcRewriteRule
}

func NewRewriter() *SvcRewriter {
	onceRewriter.Do(func() {
		svcRewriter = &SvcRewriter{
			mutex:    sync.RWMutex{},
			svcRules: make([]*SvcRewriteRule, 0, 10),
			stopC:    make(chan struct{}),
			ruleC:    make(chan *SvcRewriteRule),
		}
		//svcMatcher.Start()

	})
	return svcRewriter
}

func (rewriter *SvcRewriter) Rewrite(path string) (string, error) {
	rewriter.mutex.RLock()
	defer rewriter.mutex.RUnlock()
	for _, r := range rewriter.svcRules {
		return r.Rewrite(path)
	}
	return "", errors.New("service not found")
}

func (rewriter *SvcRewriter) In() chan<- *SvcRewriteRule {
	return rewriter.ruleC
}

func (rewriter *SvcRewriter) delete(rule *SvcRewriteRule) {
	for idx, r := range rewriter.svcRules {
		if r.Is(rule) {
			rewriter.svcRules = append(rewriter.svcRules[0:idx], rewriter.svcRules[idx+1:]...)
		}
	}
}
func (rewriter *SvcRewriter) update(rule *SvcRewriteRule) {
	rewriter.mutex.Lock()
	defer rewriter.mutex.Unlock()
	config.Logger.Info("update svc rewriter", zap.Any("rule", rule))
	switch rule.EventType {
	case dao.EventDelete:
		rewriter.delete(rule)
	case dao.EventUpdate:
		rewriter.delete(rule)
		rewriter.svcRules = append(rewriter.svcRules, rule)
	case dao.EventCreate:
		rewriter.svcRules = append(rewriter.svcRules, rule)
	default:

	}
}

func (rewriter *SvcRewriter) runLoop() {
loop:
	for {
		select {
		case rule, ok := <-rewriter.ruleC:
			if !ok {
				break loop
			}
			if err := rule.genPattern(); err != nil {
				config.Logger.Error("generate pattern", zap.Error(err))
			} else {
				rewriter.update(rule)
			}
		case <-rewriter.stopC:
			break loop
		}
	}
}

func (rewriter *SvcRewriter) Start() {
	go rewriter.runLoop()
}

func (rewriter *SvcRewriter) Stop() {
	close(rewriter.stopC)
}
