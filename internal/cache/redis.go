package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/khanhnp-2797/echo-realworld-api/internal/config"
)

type redisCache struct {
	client *redis.Client
}

// NewRedisClient creates a raw *redis.Client from the given config.
// Use this when you need to share a single client across multiple subsystems
// (e.g. cache + pub/sub hub). The caller is responsible for closing the client.
func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

// NewRedisCache creates a Redis-backed Cache from the given config.
func NewRedisCache(cfg config.RedisConfig) Cache {
	return &redisCache{client: NewRedisClient(cfg)}
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil // cache miss
	}
	return val, err
}

func (r *redisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}
