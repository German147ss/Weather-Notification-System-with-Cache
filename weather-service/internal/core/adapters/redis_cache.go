package adapters

import (
	"app/internal/core/ports"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCacheService struct {
	client *redis.Client
}

func NewRedisCacheService(client *redis.Client) ports.CacheService {
	return &RedisCacheService{client: client}
}

func (r *RedisCacheService) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCacheService) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}
