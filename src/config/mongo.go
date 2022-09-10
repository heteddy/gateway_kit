// @Author : detaohe
// @File   : mongodb.go
// @Description:
// @Date   : 2022/9/4 20:55

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway_kit/util/mongodb"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoEngine *mongodb.Client

func InitMongo(c mongodb.Config, runMode string) {
	if b, err := json.MarshalIndent(c, "", "    "); err != nil {
		panic(err)
	} else {
		fmt.Println("mongodb config:", string(b))
	}
	var e error
	MongoEngine, e = mongodb.New(c)
	if e != nil {
		panic(e)
	}
	if err := MongoEngine.Client.Ping(context.Background(), readpref.Secondary()); err != nil {
		panic(err)
	}

}
