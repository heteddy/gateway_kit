// @Author : detaohe
// @File   : svr
// @Description:
// @Date   : 2022/8/30 18:04

package admin

import (
	"context"
	"gateway_kit/core/gateway"
)

type Upper interface {
	Uppercase(context.Context, string) (string, error)
}

type ServiceReg struct {
	//
	repo gateway.HttpServiceRepo //  以后想办法通过网络保持长连接或者etcd watch变化
}

// ServiceRegRequest 创建 更新 request
type ServiceRegRequest struct {
	ServiceName string
	Description string
	Host        string
	Port        string
}

func (reg *ServiceReg) CreateService() {

}
