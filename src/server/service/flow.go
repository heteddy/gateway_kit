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
	dayDao *dao.ServiceDayDao
}

func NewFlowService(svc *dao.ServiceHourDao, req *dao.RequestHourDao, day *dao.ServiceDayDao) *FlowService {
	return &FlowService{
		svcDao: svc,
		reqDao: req,
		dayDao: day,
	}
}

type FlowInfo struct {
	Name  string
	Count int64
	Date  time.Time
	Hour  int
}

func (svc *FlowService) GetServiceDetail(ctx context.Context, service string, from, end time.Time) ([]*dao.ServiceDayEntity, error) {
	return svc.dayDao.GetDetail(ctx, "service", service, from, end)
}

func (svc *FlowService) GetServicesSum(ctx context.Context, from, end time.Time) ([]*dao.ServiceSumEntity, error) {
	return svc.dayDao.GetSum(ctx, "service", "", from, end)
}

func (svc *FlowService) GetGwDetail(ctx context.Context, name string, from, end time.Time) ([]*dao.ServiceDayEntity, error) {
	return svc.dayDao.GetDetail(ctx, "gateway", "", from, end)
}

func (svc *FlowService) GetGwSum(ctx context.Context, from, end time.Time) ([]*dao.ServiceSumEntity, error) {
	return svc.dayDao.GetSum(ctx, "gateway", "", from, end)
}

func (svc *FlowService) GetReqDetail(ctx context.Context, service string, from, end time.Time) ([]*dao.ReqHourEntity, error) {
	return svc.reqDao.GetServiceRequestsDetail(ctx, service, from, end)
}

func (svc *FlowService) GetReqSum(ctx context.Context, service, uri, method string, from, end time.Time) ([]*FlowInfo, error) {
	svc.reqDao.GetReqSum(ctx, service, uri, method, from, end)
	return nil, nil
}
