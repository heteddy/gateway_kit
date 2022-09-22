// @Author : detaohe
// @File   : mongodb.go
// @Description:
// @Date   : 2022/4/17 8:21 PM

package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Client struct {
	*mongo.Client
	*mongo.Database
	config *Config
	opts   *options.ClientOptions
}

type Config struct {
	Hosts    []string `json:"hosts"`
	User     string   `json:"user"`
	Pass     string   `json:"pass"`
	Replica  string   `json:"replica"`
	Database string   `json:"database"`
}

// New 创建一个自定义的client，继承自*mongodb.Client
func New(c Config, opts ...*options.ClientOptions) (*Client, error) {
	opt := options.Client().
		SetHosts(c.Hosts).
		SetReplicaSet(c.Replica).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(20).
		SetMinPoolSize(5).
		SetReadPreference(readpref.Secondary()) //默认读从库

	newOpts := make([]*options.ClientOptions, 0, len(opts)+1)
	newOpts = append(newOpts, opt)
	newOpts = append(newOpts, opts...)

	merged := options.MergeClientOptions(newOpts...)

	client, err := mongo.Connect(context.Background(), merged)
	if err != nil {
		return nil, err
	}
	//
	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Println(err)
	} else {
		log.Println("连接成功")
	}
	db := client.Database(c.Database)
	collectionNames, err2 := db.ListCollectionNames(context.Background(), options.ListCollections())
	if err2 != nil {
		log.Println("error of collections", err)
	} else {
		log.Println("collection names=", collectionNames)
	}
	mongoClient := &Client{
		config:   &c,
		Client:   client,
		Database: db,
		opts:     merged, // 方便打印
	}
	return mongoClient, nil
}

type Dao struct {
	Client        *Client
	Table         string
	IndexParamMap map[string]mongo.IndexModel
}

func (m Dao) Collection() *mongo.Collection {
	return m.Client.Collection(m.Table)
}

func (m Dao) CreateIndex(model mongo.IndexModel) {
	indexView := m.Collection().Indexes()
	index, err := indexView.CreateOne(context.Background(), model)
	if err != nil {
		log.Fatalf("index create failure, err =%v\n", err)

	} else {
		log.Println("index created", index)
	}
}

func (m Dao) Delete(ctx context.Context, _id string) error {
	if objID, err := primitive.ObjectIDFromHex(_id); err != nil {
		return err
	} else {
		if _, err2 := m.Collection().DeleteOne(ctx, bson.M{"_id": objID}); err2 != nil {
			return err2
		}
		return nil
	}
}

func (m Dao) SoftDelete(ctx context.Context, _id string) error {
	if objID, err := primitive.ObjectIDFromHex(_id); err != nil {
		return err
	} else {
		updatedAt := time.Now().Unix()
		ret, err := m.Collection().UpdateByID(ctx, objID, bson.M{"deleted_at": updatedAt, "updated_at": updatedAt}, options.Update().SetUpsert(false))
		if err != nil {
			return err
		} else {
			fmt.Printf("delete id=%s, deleted_at=%d, updated_count=%d\n", _id, updatedAt, ret.ModifiedCount)
		}
		return nil
	}
}
