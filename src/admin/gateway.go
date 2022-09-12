// @Author : detaohe
// @File   : gateway.go
// @Description:
// @Date   : 2022/9/10 22:38

package admin

import (
	"context"
	"gateway_kit/dao"
)

type GatewaySvc struct {
	*dao.GatewayDao
}

func NewGatewaySvc() *GatewaySvc {
	return &GatewaySvc{
		GatewayDao: dao.NewGatewayDao(),
	}
}

func (gateway *GatewaySvc) Create(ctx context.Context, entity *dao.GatewayEntity) (*dao.GatewayEntity, error) {
	return gateway.GatewayDao.Insert(ctx, entity)
}

func (gateway *GatewaySvc) Update(ctx context.Context, svcID string, entity *dao.GatewayEntity) (*dao.GatewayEntity, error) {
	return gateway.GatewayDao.Update(ctx, svcID, entity)
}

func (gateway *GatewaySvc) Delete(ctx context.Context, svcID string) error {
	return gateway.GatewayDao.SoftDelete(ctx, svcID)
}
