package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/wisle25/be-template/applications/cache"
	"github.com/wisle25/be-template/commons"
	"time"
)

type RedisCache struct /* implements Cache */ {
	redis *redis.Client
}

func NewRedisCache(redis *redis.Client) cache.Cache {
	return &RedisCache{
		redis: redis,
	}
}

func (r *RedisCache) SetCache(key string, value interface{}, expiration time.Duration) {
	ctx := context.TODO()
	err := r.redis.Set(ctx, key, value, expiration).Err()

	if err != nil {
		commons.ThrowServerError("redis_cache_err: set cache:", err)
	}
}

func (r *RedisCache) GetCache(key string) interface{} {
	ctx := context.TODO()
	val, err := r.redis.Get(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		} else {
			commons.ThrowServerError("redis_cache_err: get cache:", err)
		}
	}

	return val
}

func (r *RedisCache) DeleteCache(key string) {
	ctx := context.TODO()
	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		commons.ThrowServerError("redis_cache_err: delete cache:", err)
	}
}
