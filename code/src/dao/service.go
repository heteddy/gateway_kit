// @Author : detaohe
// @File   : service.go
// @Description:
// @Date   : 2022/4/17 8:15 PM

package dao

import "time"

type ServiceEntity struct {
	ID          string `json:"id" bson:"_id"`
	Type        string `json:"type" bson:"type"` // grpc http tcp 等
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Deleted     string `json:"deleted" bson:"deleted"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ServiceDao struct {
}
