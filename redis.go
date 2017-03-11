package main

import (
	"vraIntegrationDemo/config"

	redis "gopkg.in/redis.v5"
)

func redisNewClient() *redis.Client {
	options := config.GetRedisConf()
	client := redis.NewClient(&options)
	return client
}
