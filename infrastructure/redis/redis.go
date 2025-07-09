package redis

type RedisClient struct {
	CacheClient *CacheClient
	LockClient  *LockClient
}

func Connect() *RedisClient {
	return &RedisClient{
		CacheClient: NewCacheClient(),
		LockClient:  NewLockClient(),
	}
}
