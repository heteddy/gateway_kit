// @Author : detaohe
// @File   : config.go
// @Description:
// @Date   : 2022/9/4 20:50

package config

import (
	"gateway_kit/util"
	"gateway_kit/util/mongodb"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

var (
	GitVersion   string
	GoVersion    string
	BuildTime    string
	Environments map[string]string
	All          Configure
)

type Configure struct {
	Name    string
	Mode    string // release debug
	Version string
	Server  struct {
		HttpPort string
		Timeout  int
	}
	Gateway struct {
		HttpPort string
		GrpcPort string
		Timeout  int
	}
	RateLimit struct {
		Limit int
		Burst int
	}
	MongoC mongodb.Config
}

func InitConfigure(configureFile, logPath string) error {
	initEnvironments()
	vp := viper.New()
	vp.SetConfigFile(configureFile)
	if err := vp.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read configure %v \n", err)
	}
	for k, v := range Environments {
		if err := vp.BindEnv(k, v); err != nil {
			log.Fatalf("Failed to bind environment %v \n", err)
		}
	}
	vp.AutomaticEnv()
	if err := vp.Unmarshal(&All); err != nil {
		log.Fatalf("Failed to Unmarshal configure %v \n", err)
	}

	//if b, err := json.MarshalIndent(All, "", "    "); err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("Configured", string(b))
	//}
	gin.SetMode(All.Mode)
	util.TranslateValidator()
	if loggerErr := InitLogger(logPath, All.Name); loggerErr != nil {
		panic(loggerErr)
	}
	InitMongo(All.MongoC, All.Mode)
	// 初始化
	//gin.SetMode("test")

	return nil
}
