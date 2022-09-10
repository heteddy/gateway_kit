// @Author : detaohe
// @File   : svc.go
// @Description:
// @Date   : 2022/4/17 8:15 PM

package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// ServiceEntity 写入数据库
// 当使用独立的服务发现的时候，比如etcd，loadbalancer 通过服务发现传入servicename，然后获取所有的服务地址
// 并通过服务地址获取真实服务器地址；
// 当config.yaml中配置为k8s true说明直接通过servicename访问即可
type ServiceEntity struct {
	ID          *primitive.ObjectID `json:"id" bson:"_id"`
	Type        string              `json:"type" bson:"type"` // grpc http tcp 等
	Name        string              `json:"name" bson:"name"` //
	Description string              `json:"description" bson:"description"`
	BlackList   []string            `json:"black_list"`
	WhiteList   []string            `json:"white_list"`
	ClientFlow  int                 `json:"client_flow"` // 客户端流量控制
	ServiceFlow int                 `json:"svc_flow"`    // 服务端流量控制
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   int64 // 删除时间
}

type ServiceDao struct {
}
