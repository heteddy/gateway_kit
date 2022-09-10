// @Author : detaohe
// @File   : http.go
// @Description:
// @Date   : 2022/9/6 21:34

package dao

import (
	"gateway_kit/config"
	mongo2 "gateway_kit/util/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//type SvcAddr struct {
//	Host string `json:"host" bson:"host"`
//	Port string `json:"port" bson:"port"`
//}

type HttpSvcEntity struct {
	ID             *primitive.ObjectID `json:"id" bson:"_id"`
	Name           string              `json:"name" bson:"name"` //http服务，
	Description    string              `json:"description" bson:"description"`
	Addrs          []string            `json:"addrs" bson:"addrs"` // k8s 系统才使用,其他时候从服务发现中获取
	BlackList      []string            `json:"black_list" bson:"black_list"`
	WhiteList      []string            `json:"white_list" bson:"white_list"`
	ClientFlow     int                 `json:"client_flow" bson:"client_flow"`                                           // 客户端流量控制
	SvrFlow        int                 `json:"svr_flow" bson:"svr_flow"`                                                 // 服务端流量控制
	Category       int                 `json:"category"  bson:"category" description:"匹配类型 domain=域名, url_prefix=url前缀"` // 如果gateway绑定多个域名，可以通过访问的host，来进行重定向
	IsHttps        int                 `json:"need_https" bson:"is_https" description:"type=支持https 1=支持"`
	IsWebsocket    int                 `json:"need_websocket" bson:"is_websocket" description:"启用websocket 1=启用"`
	StripUri       []string            `json:"strip_uri" bson:"strip_uri" description:"启用strip_uri, 去掉的uri前缀"` // 如果修改url可以通过gateway修改
	UrlRewrite     []string            `json:"url_rewrite" bson:"url_rewrite" description:"url重写功能，每行一个"`      // todo 需要支持正则表达式？,当修改了uri，可以对客户端保持兼容
	HeaderTransfer []string            `json:"header_transfer" bson:"header_transfer"  description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue	"`
	CreatedAt      time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at" bson:"updated_at"`
	DeletedAt      int64               `json:"deleted_at" bson:"deleted_at" description:"deleted"` // 删除时间
}

type HttpSvcDao struct {
	mongo2.Dao
}

func NewHttpSvcDao() *HttpSvcDao {
	indices := make(map[string]mongo.IndexModel)
	idxSvcName := "idx_svc_name_deleted"
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
	return &HttpSvcDao{
		Dao: mongo2.Dao{
			Client:        config.MongoEngine,
			Table:         "",
			IndexParamMap: indices,
		},
	}
}
