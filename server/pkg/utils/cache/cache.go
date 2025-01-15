package cache

import (
	"errors"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error // 设置缓存
	Get(key string) (interface{}, bool)                                // 获取缓存
	Delete(key string) error                                           // 删除缓存
	Flush() error                                                      // 清空缓存
	Increment(key string, value int64) error                           // 增加缓存
	Decrement(key string, value int64) error                           // 减少缓存
}

// CacheType 定义缓存类型
type CacheType string

const (
	CacheTypeLocal CacheType = "local"
	CacheTypeRedis CacheType = "redis"
)

// NewCache 创建缓存实例
func NewCache(cacheType CacheType, redisAddr, redisPassword string, redisDB int) (Cache, error) {
	switch cacheType {
	case CacheTypeLocal:
		return NewLocalCache(), nil
	case CacheTypeRedis:
		return NewRedisCache(redisAddr, redisPassword, redisDB), nil
	default:
		return nil, errors.New("unsupported cache type")
	}
}
