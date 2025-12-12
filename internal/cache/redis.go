package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisClient{Client: rdb}
}
