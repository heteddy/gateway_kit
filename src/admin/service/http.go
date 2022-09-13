// @Author : detaohe
// @File   : server.go
// @Description:
// @Date   : 2022/9/10 22:38

package service

import (
	"context"
	"gateway_kit/dao"
)

type ServiceSvr struct {
	*dao.HttpSvcDao
}

func NewServiceSvr() *ServiceSvr {
	svr := &ServiceSvr{
		HttpSvcDao: dao.NewHttpSvcDao(),
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

func (svr *ServiceSvr) GetServiceByName(ctx context.Context, name string) ([]*dao.HttpSvcEntity, error) {
	return svr.HttpSvcDao.GetSvc(ctx, name)
}

func (svr *ServiceSvr) GetServiceByID(ctx context.Context, id string) (*dao.HttpSvcEntity, error) {
	return svr.HttpSvcDao.GetByID(ctx, id)
}
