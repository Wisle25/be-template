package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/wisle25/be-template/applications/database"
	"time"
)

// RedisCache implements Cache
type RedisCache struct {
	redis *redis.Client
}

func NewRedisCache(redis *redis.Client) database.Cache {
	return &RedisCache{
		redis: redis,
	}
}

func (r *RedisCache) SetCache(key string, value interface{}, expiration time.Duration) {
	ctx := context.TODO()
	err := r.redis.Set(ctx, key, value, expiration).Err()

	if err != nil {
		panic(fmt.Errorf("redis_cache_err: %v", err))
	}
}

func (r *RedisCache) GetCache(key string) interface{} {
	ctx := context.TODO()
	val, err := r.redis.Get(ctx, key).Result()

	if err != nil {
		panic(fmt.Errorf("redis_cache_err: %v", err))
	}

	return val
}
