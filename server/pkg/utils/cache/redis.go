package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache 实现 Cache 接口
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建一个新的 RedisCache
func NewRedisCache(addr, password string, db int) *RedisCache {
	// edis.NewClusterClient(&redis.ClusterOptions{ // 集群集群
	// Addrs:    ClusterAddrs,
	// Password: Password,
	// })
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{
		client: client,
	}
}

// Set 设置缓存
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func (r *RedisCache) Get(key string) (interface{}, bool) {
	ctx := context.Background()
	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false
	} else if err != nil {
		return nil, false
	}
	return value, true
}

// Delete 删除缓存
func (r *RedisCache) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

// Flush 清空缓存
func (r *RedisCache) Flush() error {
	ctx := context.Background()
	return r.client.FlushAll(ctx).Err()
}

// Increment 增加缓存
func (r *RedisCache) Increment(key string, value int64) error {
	ctx := context.Background()
	return r.client.IncrBy(ctx, key, value).Err()
}

// Decrement 减少缓存
func (r *RedisCache) Decrement(key string, value int64) error {
	ctx := context.Background()
	return r.client.DecrBy(ctx, key, value).Err()
}
