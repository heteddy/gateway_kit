// @Author : detaohe
// @File   : redis
// @Description:
// @Date   : 2022/10/15 16:02

package config

import (
	"fmt"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func InitRedis(addr, pass string, db int) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
		OnConnect: func(conn *redis.Conn) error {
			ret, err := conn.Ping().Result()
			if err != nil {
				panic(err)
			}
			fmt.Printf("redis connection:%s\n", ret)
			return nil
		},
	})
	if RedisClient == nil {
		panic("connect redis error")
	}
	ret2, err2 := RedisClient.Ping().Result()
	if err2 != nil {
		panic(err2)
	}
	fmt.Printf("redis client:%s\n", ret2)
}
