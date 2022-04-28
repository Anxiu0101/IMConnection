package cache

import (
	"IMConnection/conf"
	"IMConnection/pkg/logging"
	"context"
	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
	Ctx         context.Context
)

func Setup() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:        conf.RedisSetting.Host,
		Password:    conf.RedisSetting.Password,
		DB:          0,
		IdleTimeout: conf.RedisSetting.IdleTimeout,
	})

	Ctx = context.Background()

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		logging.Info("redis fail to connect", err.Error())
		return
	}
}
