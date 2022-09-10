// @Author : detaohe
// @File   : svr.go
// @Description:
// @Date   : 2022/4/17 8:15 PM

package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// GatewayEntity 写入数据库
// 当使用独立的服务发现的时候，比如etcd，loadbalancer 通过服务发现传入servicename，然后获取所有的服务地址
// 并通过服务地址获取真实服务器地址；
// 当config.yaml中配置为k8s true说明直接通过servicename访问即可
type GatewayEntity struct {
	ID          *primitive.ObjectID `json:"id" bson:"_id"`
	Name        string              `json:"name" bson:"name"`               // gateway name
	Description string              `json:"description" bson:"description"` //描述
	BlockList   []string            `json:"block_list"`                     // 网关黑名单，所有的服务通用
	AllowList   []string            `json:"allow_list"`                     // 网关黑名单，所有的服务通用
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   int64 // 删除时间
}

/*

黑白名单过滤流程：
	临时规则白名单-可以直接跳过，在context中加入一个特殊字段

	gateway 白名单pass
	gateway 黑名单直接block

	service 白名单pass
	service 黑名单直接block

*/

type GatewayDao struct {
}
