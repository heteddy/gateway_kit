// @Author : detaohe
// @File   : svc.go
// @Description:
// @Date   : 2022/4/17 8:15 PM

package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ServiceEntity struct {
	ID          *primitive.ObjectID `json:"id" bson:"_id"`
	Type        string              `json:"type" bson:"type"` // grpc http tcp 等
	Name        string              `json:"name" bson:"name"`
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
