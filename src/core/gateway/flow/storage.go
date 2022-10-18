// @Author : detaohe
// @File   : storage
// @Description:
// @Date   : 2022/10/16 18:43

package flow

import (
	"gateway_kit/config"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"strconv"
)

type FlowStorage struct {
	client      *redis.Client
	storageChan chan StorageCmd
	stopC       chan struct{}
}

func NewFlowStorage() *FlowStorage {
	return &FlowStorage{
		client:      config.RedisClient,
		storageChan: make(chan StorageCmd, 50),
		stopC:       make(chan struct{}),
	}
}
func (storage *FlowStorage) loadByKey(key string) (map[string]int64, error) {
	ret, err := storage.client.Get(key).Result()
	if err != nil {
		config.Logger.Error("load redis key error", zap.Error(err))
		return nil, err
	} else {
		value, err2 := strconv.ParseInt(ret, 10, 64)
		if err2 != nil {
			config.Logger.Error("convert string value error", zap.Error(err), zap.String(key, ret))
			return nil, err
		}
		return map[string]int64{key: value}, nil
	}
}

func (storage *FlowStorage) loadByPrefix(prefix string) (map[string]int64, error) {
	keys, err := storage.client.Keys(prefix).Result()
	if err != nil {
		config.Logger.Error("load redis keys error", zap.Error(err))
		return nil, err
	} else {
		result := make(map[string]int64)
		values, err2 := storage.client.MGet(keys...).Result()
		if err2 != nil {
			return result, err2
		} else {
			for idx, v := range values {
				_v := v.(string)
				value, err3 := strconv.ParseInt(_v, 10, 64)
				if err3 != nil {
					config.Logger.Error("convert string value error", zap.Error(err3), zap.String("s", _v))
					return nil, err3
				} else {
					if value > 0 {
						result[keys[idx]] = value
					}
				}
			}
		}
		return result, nil
	}
}
func (storage *FlowStorage) runLoop() {
loop:
	for {
		select {
		case cmd, ok := <-storage.storageChan:
			if !ok {
				break loop
			}
			if err := cmd.Run(storage); err != nil {
				config.Logger.Error("run command error", zap.Error(err))
			}

		case <-storage.stopC:
			break loop
		}
	}
}

func (storage *FlowStorage) Stop() {
	close(storage.stopC)
	close(storage.storageChan)
}

func (storage *FlowStorage) Start() {
	go storage.runLoop()
}

func (storage *FlowStorage) In() chan<- StorageCmd {
	return storage.storageChan
}

func (storage *FlowStorage) DecrBy(key string, value int64) error {
	return storage.client.DecrBy(key, value).Err()
}

func (storage *FlowStorage) IncrBy(key string, value int64) error {
	return storage.client.IncrBy(key, value).Err()
}
