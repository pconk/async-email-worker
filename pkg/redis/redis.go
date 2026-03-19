package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient(redisAddress string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
}
