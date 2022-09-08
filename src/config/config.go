// @Author : detaohe
// @File   : config.go
// @Description:
// @Date   : 2022/9/4 20:50

package config

import (
	"encoding/json"
	"fmt"
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
	Service string
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
}

func InitConfigure(configureFile, logPath string) error {
	//All.Service = "teddy-gateway"
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
	if b, err := json.MarshalIndent(All, "", "    "); err != nil {
		panic(err)
	} else {
		fmt.Println("Configured", string(b))
	}

	if loggerErr := InitLogger(logPath, All.Service); loggerErr != nil {
		panic(loggerErr)
	}
	gin.SetMode(All.Mode)

	return nil
}
