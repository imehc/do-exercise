package cache

import (
	"time"

	"github.com/songzhibin97/gkit/cache/local_cache"
)

// LocalCache 实现 Cache 接口
type LocalCache struct {
	cache *local_cache.Cache
}

// NewLocalCache 创建一个新的 LocalCache
func NewLocalCache() *LocalCache {
	return &LocalCache{
		cache: &local_cache.Cache{},
	}
}

// Set 设置缓存
func (l *LocalCache) Set(key string, value interface{}, expiration time.Duration) error {
	l.cache.Set(key, value, expiration)
	return nil
}

// Get 获取缓存
func (l *LocalCache) Get(key string) (interface{}, bool) {
	return l.cache.Get(key)
}

// Delete 删除缓存
func (l *LocalCache) Delete(key string) error {
	l.cache.Delete(key)
	return nil
}

func (l *LocalCache) Flush() error {
	l.cache.Flush()
	return nil
}

func (l *LocalCache) Increment(key string, value int64) error {
	l.cache.Increment(key, value)
	return nil
}

func (l *LocalCache) Decrement(key string, value int64) error {
	l.cache.Decrement(key, value)
	return nil
}
