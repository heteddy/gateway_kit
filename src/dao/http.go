// @Author : detaohe
// @File   : server.go
// @Description:
// @Date   : 2022/9/6 21:34

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

//type SvcAddr struct {
//	Host string `json:"host" bson:"host"`
//	Port string `json:"port" bson:"port"`
//}

const (
	SvcCategoryUrlPrefix = 0
	SvcCategoryHost      = 1
)

type HttpSvcEntity struct {
	ID             *primitive.ObjectID `json:"id" bson:"_id"`
	Name           string              `json:"name" bson:"name"` //http服务，
	Description    string              `json:"description" bson:"description"`
	Addrs          []string            `json:"addrs" bson:"addrs"` // k8s 系统才使用,其他时候从服务发现中获取
	BlockList      []string            `json:"block_list" bson:"block_list"`
	AllowList      []string            `json:"allow_list" bson:"allow_list"`
	ClientQps      int                 `json:"client_qps" bson:"client_qps"`                                             // 客户端流量控制
	ServerQps      int                 `json:"server_qps" bson:"server_qps"`                                             // 服务端流量控制
	Category       int                 `json:"category"  bson:"category" description:"匹配类型 domain=域名, url_prefix=url前缀"` // 如果gateway绑定多个域名，可以通过访问的host，来进行重定向
	MatchRule      string              `json:"match_rule" bson:"match_rule"`                                             // 匹配的项目与category结合使用，如果是domain，host==match_rule，否则是url前缀匹配
	IsHttps        bool                `json:"need_https" bson:"is_https" description:"type=支持https 1=支持"`
	IsWebsocket    bool                `json:"need_websocket" bson:"is_websocket" description:"启用websocket 1=启用"`
	StripUri       []string            `json:"strip_uri" bson:"strip_uri" description:"启用strip_uri, 去掉的uri前缀"` // 如果修改url可以通过gateway修改
	UrlRewrite     []string            `json:"url_rewrite" bson:"url_rewrite" description:"url重写功能，每行一个"`      // todo 需要支持正则表达式？,当修改了uri，可以对客户端保持兼容
	HeaderTransfer []string            `json:"header_transfer" bson:"header_transfer"  description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue	"`
	CreatedAt      time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at" bson:"updated_at"`
	DeletedAt      int64               `json:"deleted_at" bson:"deleted_at" description:"deleted"` // 删除时间
}

type HttpSvcDao struct {
	mongodb.Dao
}

func NewHttpSvcDao() *HttpSvcDao {
	indices := make(map[string]mongo.IndexModel)
	idxSvcName := "idx_svc_name_category_rule_deleted"
	indexBackground := true
	unique := true

	indices[idxSvcName] = mongo.IndexModel{
		Keys: bson.D{{"name", 1}, {"deleted_at", 1}, {"category", 1}, {"rule", 1}},
		Options: &options.IndexOptions{
			Name:       &idxSvcName,
			Background: &indexBackground,
			Unique:     &unique,
		},
	}
	return &HttpSvcDao{
		Dao: mongodb.Dao{
			Client:        config.MongoEngine,
			Table:         "",
			IndexParamMap: indices,
		},
	}
}

func (engine *HttpSvcDao) All(ctx context.Context) (entities []*HttpSvcEntity, err error) {
	opt := options.Find().SetSort(bson.D{{"updated_at", -1}})
	var cursor *mongo.Cursor
	cursor, err = engine.Collection().Find(ctx, bson.M{"deleted_at": 0}, opt)
	if err != nil {
		return
	} else {
		err = cursor.All(ctx, &entities)
		return
	}
}

func (engine *HttpSvcDao) GetSvc(ctx context.Context, svc string) (entities []*HttpSvcEntity, err error) {
	opt := options.Find().SetSort(bson.D{{"updated_at", -1}})
	var cursor *mongo.Cursor
	cursor, err = engine.Collection().Find(ctx, bson.D{{"name", svc}, {"deleted_at", 0}}, opt)
	if err != nil {
		return
	} else {
		err = cursor.All(ctx, &entities)
		return
	}
}

func (engine *HttpSvcDao) GetByID(ctx context.Context, id string) (entity *HttpSvcEntity, err error) {
	var objID primitive.ObjectID
	if objID, err = primitive.ObjectIDFromHex(id); err != nil {
		return nil, err
	} else {
		err = engine.Collection().FindOne(ctx, bson.M{"_id": objID}).Decode(&entity)
		if err != nil {
			return nil, err
		} else {
			return
		}

	}
}

func (engine *HttpSvcDao) Delete(ctx context.Context, _id string) error {
	if objID, err := primitive.ObjectIDFromHex(_id); err != nil {
		return err
	} else {
		if _, err2 := engine.Collection().DeleteOne(ctx, bson.M{"_id": objID}); err2 != nil {
			return err2
		}
		return nil
	}
}

func (engine *HttpSvcDao) Insert(ctx context.Context, svc *HttpSvcEntity) (*HttpSvcEntity, error) {
	svc.CreatedAt = time.Now()
	svc.UpdatedAt = time.Now()
	result, err := engine.Collection().InsertOne(ctx, svc)
	if err != nil {
		return svc, err
	} else {
		if objID, ok := result.InsertedID.(primitive.ObjectID); ok {
			svc.ID = &objID
			return svc, nil
		}
		return svc, errors.New("mongo _id类型转换错误")
	}
}

func (engine *HttpSvcDao) Update(ctx context.Context, _id string, svc *HttpSvcEntity) (*HttpSvcEntity, error) {
	//service.CreatedAt = time.Now()
	svc.UpdatedAt = time.Now()
	if objID, err := primitive.ObjectIDFromHex(_id); err != nil {
		return svc, err
	} else {

		ret, err := engine.Collection().UpdateByID(ctx, objID, svc, options.Update().SetUpsert(false))
		if err != nil {
			return svc, err
		} else {
			config.Logger.Info("update ", zap.String("_id", _id), zap.Int64("count", ret.ModifiedCount))
			return svc, nil
		}

	}
}
