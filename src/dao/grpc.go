// @Author : detaohe
// @File   : grpc.go
// @Description:
// @Date   : 2022/9/6 21:41

package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GrpcSvcEntity struct {
	ID             *primitive.ObjectID `json:"id" `
	Name           string
	Description    string
	Addrs          []string  `json:"addrs"`
	HeaderTransfer []string  `json:"header_transfer" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      int64     `json:"deleted_at"` // 删除时间
}
