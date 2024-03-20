package infrastructure

import (
	"context"

	"github.com/go-redis/redis"
)

func NewRedisClient(ctx context.Context, config Config) redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return *rdb
}

type Cache interface {
	SetObj(ctx context.Context, key string, obj interface{}, ttlInSec int) (err error)
	GetObj(ctx context.Context, key string) (result *interface{}, err error)
	Del(ctx context.Context, key string) (err error)
}

type redisCache struct {
	client redis.Client
}

func NewCache(redisClient redis.Client) Cache {
	return &redisCache{
		client: redisClient,
	}
}

func (cache *redisCache) SetObj(ctx context.Context, key string, obj interface{}, ttlInSec int) (err error) {

	return
}

func (cache *redisCache) GetObj(ctx context.Context, key string) (result *interface{}, err error) {

	return
}

func (cache *redisCache) Del(ctx context.Context, key string) (err error) {

	return
}
