// @Author : detaohe
// @File   : mongodb.go
// @Description:
// @Date   : 2022/4/17 8:21 PM

package mongodb

import (
	"context"
	"fmt"
	"gateway_kit/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
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
	toDelIndex    []string
}

func (d Dao) Collection() *mongo.Collection {
	return d.Client.Collection(d.Table)
}

func (d Dao) CreateIndexes() {
	c, err := d.Collection().Indexes().List(context.Background())
	if err != nil {
		config.Logger.Error("err of list index", zap.Error(err))
	} else {
		var indexes []bson.M
		if err := c.All(context.Background(), &indexes); err != nil {
			config.Logger.Error("err of all index", zap.Error(err))
		} else {
			istatus := make(map[string]bool, len(d.IndexParamMap))
			for idx, index := range indexes {
				for k, v := range index {
					if k == "name" {
						idxName, ok := v.(string)
						if !ok {
							continue
						}
						istatus[idxName] = true
						config.Logger.Info("已存在 index", zap.Int("idx", idx), zap.String("name", idxName))
					}

				}
			}
			for _, name := range d.toDelIndex {
				indexView := d.Collection().Indexes()
				_, err := indexView.DropOne(context.Background(), name)
				if err != nil {
					config.Logger.Error("err of del index", zap.Error(err), zap.String("name", name))
				} else {
					config.Logger.Info("success of del index", zap.String("name", name))
				}
			}
			for k, _ := range d.IndexParamMap {
				if _, exists := istatus[k]; !exists {
					d.CreateIndex(d.IndexParamMap[k])
				}
			}
		}
	}
}

func (d Dao) CreateIndex(model mongo.IndexModel) {
	indexView := d.Collection().Indexes()
	index, err := indexView.CreateOne(context.Background(), model)
	if err != nil {
		log.Fatalf("index create failure, err =%v\n", err)

	} else {
		log.Println("index created", index)
	}
}

func (d Dao) Delete(ctx context.Context, _id string) error {
	if objID, err := primitive.ObjectIDFromHex(_id); err != nil {
		return err
	} else {
		if _, err2 := d.Collection().DeleteOne(ctx, bson.M{"_id": objID}); err2 != nil {
			return err2
		}
		return nil
	}
}

func (d Dao) SoftDelete(ctx context.Context, _id string) error {
	if objID, err := primitive.ObjectIDFromHex(_id); err != nil {
		return err
	} else {
		now := time.Now()
		updatedAt := now.Unix()
		ret, err := d.Collection().UpdateByID(ctx, objID, bson.M{"$set": bson.M{"deleted_at": updatedAt, "updated_at": now}},
			options.Update().SetUpsert(false))
		if err != nil {
			return err
		} else {
			fmt.Printf("delete id=%s, deleted_at=%d, updated_count=%d\n", _id, updatedAt, ret.ModifiedCount)
		}
		return nil
	}
}
