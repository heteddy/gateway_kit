// @Author : detaohe
// @File   : client
// @Description:
// @Date   : 2022/9/9 20:58

package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TemporaryClient struct {
	ID          *primitive.ObjectID `json:"id" bson:"_id"`
	Name        string              `json:"name" bson:"name"` //http服务，
	Description string              `json:"description" bson:"description"`
	Host        string              `json:"host" bson:"host"`
	Flow        int                 `json:"flow" bson:"flow"`
	ExpiredAt   time.Time           `json:"expired_at" bson:"expired_at"`
}
