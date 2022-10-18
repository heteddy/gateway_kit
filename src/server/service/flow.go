// @Author : detaohe
// @File   : flow.go
// @Description:
// @Date   : 2022/10/18 10:48

package service

import (
	"context"
	"gateway_kit/dao"
	"time"
)

type FlowService struct {
	svcDao *dao.ServiceHourDao
	reqDao *dao.RequestHourDao
}

func NewFlowService(svc *dao.ServiceHourDao, req *dao.RequestHourDao) *FlowService {
	return &FlowService{
		svcDao: svc,
		reqDao: req,
	}
}

type FlowInfo struct {
	Name  string
	Count int64
	Date  time.Time
	Hour  int
}

func (svc *FlowService) GetServiceDetail(ctx context.Context, service string, from, end time.Time) ([]*dao.ServiceHourEntity, error) {
	return svc.svcDao.GetServicesDetail(ctx, service, from, end)
}

func (svc *FlowService) GetServicesSum(ctx context.Context, from, end time.Time) ([]*FlowInfo, error) {
	svc.svcDao.GetSum(ctx, from, end)
	return nil, nil
}

func (svc *FlowService) GetGwDetail(ctx context.Context, from, end time.Time) ([]*FlowInfo, error) {

	return nil, nil
}

func (svc *FlowService) GetGwSum(ctx context.Context, from, end time.Time) (*FlowInfo, error) {
	return nil, nil
}

func (svc *FlowService) GetReqDetail(ctx context.Context, service string, from, end time.Time) ([]*FlowInfo, error) {
	svc.reqDao.GetServiceRequestsDetail(ctx, service, from, end)
	return nil, nil
}

func (svc *FlowService) GetReqSum(ctx context.Context, service, uri, method string, from, end time.Time) ([]*FlowInfo, error) {
	svc.reqDao.GetReqSum(ctx, service, uri, method, from, end)
	return nil, nil
}
