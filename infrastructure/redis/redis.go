package redis

import "github.com/go-redis/redis/v8"

type RedisClient struct {
	CacheClient *redis.Client
	LockClient  *redis.Client
}

func Connect() *RedisClient {
	cacheClient := NewCacheClient()
	lockClient := NewLockClient()

	return &RedisClient{
		CacheClient: cacheClient.client,
		LockClient:  lockClient.client,
	}
}
