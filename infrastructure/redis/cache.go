package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
)

type CacheClient struct {
	client *redis.Client
}

func NewCacheClient() *CacheClient {
	cfg := config.AppConfig

	return &CacheClient{
		client: redis.NewClient(
			&redis.Options{
				Addr:     cfg.RedisAddress,
				Password: cfg.RedisPassword,
				DB:       0,
			},
		),
	}
}

func (c *CacheClient) Get(ctx context.Context, key string, dest any) (data string, err error) {
	data, err = c.client.Get(ctx, key).Result()

	return
}

func (c *CacheClient) Exist(ctx context.Context, key string) (exist bool, err error) {
	data, err := c.client.Exists(ctx, key).Result()

	if data == 1 {
		exist = true
	}

	return exist, err
}

func (c *CacheClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *CacheClient) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *CacheClient) Shutdown() error {
	return c.client.Close()
}
