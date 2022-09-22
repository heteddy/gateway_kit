// @Author : detaohe
// @File   : svc_name.go
// @Description:
// @Date   : 2022/4/17 8:15 PM

package dao

import (
	"context"
	"errors"
	"gateway_kit/config"
	"gateway_kit/util/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

// GatewayEntity 写入数据库
// 当使用独立的服务发现的时候，比如etcd，loadbalancer 通过服务发现传入servicename，然后获取所有的服务地址
// 并通过服务地址获取真实服务器地址；
// 当config.yaml中配置为k8s true说明直接通过servicename访问即可
type GatewayEntity struct {
	ID          *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string              `json:"name" bson:"name"`               // gateway name
	Description string              `json:"description" bson:"description"` //描述
	BlockList   []string            `json:"block_list" bson:"block_list"`   // 网关黑名单，所有的服务通用
	AllowList   []string            `json:"allow_list" bson:"allow_list"`   // 网关黑名单，所有的服务通用
	CreatedAt   time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" bson:"updated_at"`
	DeletedAt   int64               `json:"deleted_at" bson:"deleted_at"` // 删除时间
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
	mongodb.Dao
}

func NewGatewayDao() *GatewayDao {
	indices := make(map[string]mongo.IndexModel)
	idxSvcName := "idx_gw_name"
	indexBackground := true
	unique := true

	indices[idxSvcName] = mongo.IndexModel{
		Keys: bson.D{{"name", 1}, {"deleted_at", 1}},
		Options: &options.IndexOptions{
			Name:       &idxSvcName,
			Background: &indexBackground,
			Unique:     &unique,
		},
	}
	return &GatewayDao{
		Dao: mongodb.Dao{
			Client:        config.MongoEngine,
			Table:         config.All.Tables.Gateway,
			IndexParamMap: indices,
		},
	}
}

func (engine *GatewayDao) GetByName(ctx context.Context, name string) (entities []*GatewayEntity, err error) {
	opt := options.Find().SetSort(bson.D{{"updated_at", -1}})
	var cursor *mongo.Cursor
	cursor, err = engine.Collection().Find(ctx, bson.D{{"name", name}, {"deleted_at", 0}}, opt)
	if err != nil {
		return
	} else {
		err = cursor.All(ctx, &entities)
		return
	}
}

//func (engine *GatewayDao) Delete(ctx context.Context, _id string) error {
//	if objID, err := primitive.ObjectIDFromHex(_id); err != nil {
//		return err
//	} else {
//		if _, err2 := engine.Collection().DeleteOne(ctx, bson.M{"_id": objID}); err2 != nil {
//			return err2
//		}
//		return nil
//	}
//}

func (engine *GatewayDao) Insert(ctx context.Context, gw *GatewayEntity) (*GatewayEntity, error) {
	gw.CreatedAt = time.Now()
	gw.UpdatedAt = time.Now()
	result, err := engine.Collection().InsertOne(ctx, gw)
	if err != nil {
		return gw, err
	} else {
		if objID, ok := result.InsertedID.(primitive.ObjectID); ok {
			gw.ID = &objID
			return gw, nil
		}
		return gw, errors.New("mongo _id类型转换错误")
	}
}

func (engine *GatewayDao) Update(ctx context.Context, _id string, gw *GatewayEntity) (*GatewayEntity, error) {
	//gw.CreatedAt = time.Now()
	gw.UpdatedAt = time.Now()
	if objID, err := primitive.ObjectIDFromHex(_id); err != nil {
		return gw, err
	} else {

		ret, err := engine.Collection().UpdateByID(ctx, objID, gw, options.Update().SetUpsert(false))
		if err != nil {
			return gw, err
		} else {
			config.Logger.Info("update ", zap.String("_id", _id), zap.Int64("count", ret.ModifiedCount))
			return gw, nil
		}

	}
}

func (engine *GatewayDao) All(ctx context.Context) (entities []*GatewayEntity, err error) {
	var cursor *mongo.Cursor
	cursor, err = engine.Collection().Find(ctx, bson.D{{"deleted_at", 0}})
	if err != nil {
		return
	} else {
		err = cursor.All(ctx, &entities)
		return
	}
}

func (engine *GatewayDao) ChangeFrom(ctx context.Context, t time.Time) (entities []*GatewayEntity, err error) {
	opt := options.Find().SetSort(bson.D{{"updated_at", -1}})
	var cursor *mongo.Cursor
	cursor, err = engine.Collection().Find(ctx, bson.D{{"updated_at", bson.M{"$gt": t}}}, opt)
	if err != nil {
		return
	} else {
		err = cursor.All(ctx, &entities)
		return
	}
}
