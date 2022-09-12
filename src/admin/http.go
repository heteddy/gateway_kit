// @Author : detaohe
// @File   : http.go
// @Description:
// @Date   : 2022/9/10 22:38

package admin

import (
	"context"
	"gateway_kit/config"
	"gateway_kit/dao"
	"go.uber.org/zap"
	"strings"
	"sync"
)

type ServiceSvr struct {
	*dao.HttpSvcDao
	entities     []*dao.HttpSvcEntity
	svcEntityMap map[string]*dao.HttpSvcEntity // 以后放到redis，目前
	mutex        sync.RWMutex
}

func NewServiceSvr(receiver chan []*dao.HttpSvcEntity) *ServiceSvr {
	svr := &ServiceSvr{
		HttpSvcDao:   dao.NewHttpSvcDao(),
		entities:     make([]*dao.HttpSvcEntity, 0, 10),
		svcEntityMap: make(map[string]*dao.HttpSvcEntity),
		mutex:        sync.RWMutex{},
	}
	return svr
}

func (svr *ServiceSvr) Create(ctx context.Context, service *dao.HttpSvcEntity) (*dao.HttpSvcEntity, error) {
	return svr.HttpSvcDao.Insert(ctx, service)
}

func (svr *ServiceSvr) Update(ctx context.Context, sID string, service *dao.HttpSvcEntity) (*dao.HttpSvcEntity, error) {
	return svr.HttpSvcDao.Update(ctx, sID, service)
}

func (svr *ServiceSvr) Delete(ctx context.Context, sID string) error {
	return svr.HttpSvcDao.SoftDelete(ctx, sID)
}

func (svr *ServiceSvr) All(ctx context.Context) ([]*dao.HttpSvcEntity, error) {
	return svr.HttpSvcDao.All(ctx)
}

func (svr *ServiceSvr) GetService(ctx context.Context, name string) ([]*dao.HttpSvcEntity, error) {
	return svr.HttpSvcDao.GetSvc(ctx, name)
}

func (svr *ServiceSvr) GetServiceName(host, path string) (string, error) {
	svr.mutex.RLock()
	defer svr.mutex.RUnlock()
	for _, entity := range svr.entities {
		switch entity.Category {
		case dao.SvcCategoryUrlPrefix:
			// todo 增加正则表达式
			if strings.HasPrefix(path, entity.MatchRule) {
				return entity.Name, nil
			}
		case dao.SvcCategoryHost:
			if host == entity.MatchRule {
				return entity.Name, nil
			}
		default:
			//
			config.Logger.Warn(
				"error of entity category",
				zap.String("entity_id", entity.ID.Hex()),
				zap.String("entity_name", entity.Name))
		}
	}
	return "", nil
}
