package common

import (
	"context"
	"errors"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 登录失败计数相关常量与辅助方法，用于防爆破
const loginFailPrefix = "login_fail:"

// loginAttempts 返回失败阈值与锁定时长（配置为 0 时使用默认值）
func loginAttempts() (int, time.Duration) {
	max := global.Config.Auth.LoginMaxAttempts
	if max <= 0 {
		max = 5
	}
	lock := global.Config.Auth.LoginLockMinutes
	if lock <= 0 {
		lock = 5
	}
	return max, time.Duration(lock) * time.Minute
}

// loginFailKeys 生成 IP 与用户名两维度的计数 key
func loginFailKeys(ip, username string) []string {
	keys := []string{loginFailPrefix + "ip:" + ip}
	if username != "" {
		keys = append(keys, loginFailPrefix+"user:"+username)
	}
	return keys
}

// loginPenalty 返回当前 IP/用户名应承受的惩罚：是否已达硬锁阈值，以及本次尝试需等待的渐进延迟。
//
// 在到达阈值前，每多一次失败，下一次尝试的等待时间按 2 的幂递增并封顶，
// 让持续爆破越来越慢；达到阈值后直接拒绝（登录锁定）。
func loginPenalty(ip, username string) (locked bool, delay time.Duration) {
	ctx := context.Background()
	max, _ := loginAttempts()
	count := 0
	for _, key := range loginFailKeys(ip, username) {
		n, err := global.Redis.Get(ctx, key).Int()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			global.Log.Error("登录失败计数读取异常，按锁定处理",
				zap.String("key", key), zap.Error(err))
			return true, 0
		}
		if n > count {
			count = n
		}
	}
	if count >= max {
		return true, 0
	}
	return false, progressiveDelay(count)
}

// progressiveDelay 根据上一轮累计的失败次数计算本次请求需等待的延迟。
func progressiveDelay(failCount int) time.Duration {
	if failCount <= 0 {
		return 0
	}
	const (
		base    = 500 * time.Millisecond // 首次失败后的等待
		maxWait = 20 * time.Second       // 渐进延迟封顶
	)
	d := base
	for i := 1; i < failCount && d < maxWait; i++ {
		d *= 2
	}
	return d
}

// registerLoginFailure 记录一次登录失败。
// INCR 与 EXPIRE 通过管道一次发出，避免两次往返之间的竞态导致计数永不过期。
func registerLoginFailure(ip, username string) {
	ctx := context.Background()
	_, lock := loginAttempts()
	for _, key := range loginFailKeys(ip, username) {
		pipe := global.Redis.Pipeline()
		incr := pipe.Incr(ctx, key)
		// NX 保证只在没有 TTL 时设置，不会因后续失败而不断续期，
		// 否则持续攻击可以把锁定窗口无限延长。
		pipe.ExpireNX(ctx, key, lock)
		if _, err := pipe.Exec(ctx); err != nil {
			global.Log.Error("登录失败计数写入异常",
				zap.String("key", key), zap.Error(err))
			continue
		}
		_ = incr.Val()
	}
}

// clearLoginFailures 登录成功后清除失败计数
func clearLoginFailures(ip, username string) {
	ctx := context.Background()
	for _, key := range loginFailKeys(ip, username) {
		global.Redis.Del(ctx, key)
	}
}
