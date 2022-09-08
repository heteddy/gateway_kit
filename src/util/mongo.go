// @Author : detaohe
// @File   : mongo.go
// @Description:
// @Date   : 2022/4/17 8:21 PM

package util

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type MongoClient struct {
	*mongo.Client
	*mongo.Database
	config *MongoConfig
	opts   *options.ClientOptions
}

type MongoConfig struct {
	Hosts      []string `json:"hosts"`
	ReplicaSet string   `json:"replica_set"`
	Database   string   `json:"database"`
}

// New 创建一个自定义的client，继承自*mongo.Client
// @Description: 只支持副本集模式
// @param config 基本配置
// @param opts
// @return *MongoClient
// @return error
//
func New(config MongoConfig, opts ...*options.ClientOptions) (*MongoClient, error) {
	opt := options.Client().
		SetHosts(config.Hosts).
		SetReplicaSet("rs0").
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
	db := client.Database(config.Database)
	collectionNames, err2 := db.ListCollectionNames(context.Background(), options.ListCollections())
	if err2 != nil {
		log.Println("error of collections", err)
	} else {
		log.Println("collection names=", collectionNames)
	}
	mongoClient := &MongoClient{
		config:   &config,
		Client:   client,
		Database: db,
		opts:     merged, // 方便打印
	}
	return mongoClient, nil
}
