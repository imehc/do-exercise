package util

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const captchaPrefix = "captcha:"

// 验证码
type Captcha struct {
	Redis     *redis.Client
	CacheTime time.Duration
	Log       *zap.Logger
}

func NewCaptcha(redis *redis.Client, cacheTime time.Duration) *Captcha {
	if redis == nil {
		panic("Redis client is nil")
	}
	return &Captcha{
		Redis:     redis,
		CacheTime: cacheTime,
	}
}

func (c *Captcha) Set(key, value string) (err error) {
	context := context.Background()
	key = captchaPrefix + key

	err = c.Redis.Set(
		context,
		key,
		value,
		c.CacheTime,
	).Err()
	if err != nil {
		return
	}
	return nil
}

func (c *Captcha) Get(key string, clear bool) (value string) {
	context := context.Background()
	key = captchaPrefix + key
	value, err := c.Redis.Get(context, key).Result()
	if err != nil {
		c.Log.Error("redis get failed", zap.Error(err))
		return ""
	}
	if clear {
		err := c.Redis.Del(context, key).Err()
		if err != nil {
			c.Log.Error("redis del failed", zap.Error(err))
			return ""
		}
	}
	return
}

func (c *Captcha) Verify(key, value string, clear bool) bool {
	stored := c.Get(key, clear)
	if stored == "" {
		return false
	}
	return strings.EqualFold(stored, value)
}
