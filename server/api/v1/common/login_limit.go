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

// loginIsLocked 判断当前 IP/用户名是否已被锁定。
//
// Redis 异常时返回 true（失败关闭）。早期实现用 `err == nil && n >= max` 判断，
// 于是任何 Redis 超时、驱逐或宕机都会让锁定失效——恰好在监控最嘈杂、
// 最难察觉的时候放开了无限密码猜测。
func loginIsLocked(ip, username string) bool {
	ctx := context.Background()
	max, _ := loginAttempts()
	for _, key := range loginFailKeys(ip, username) {
		n, err := global.Redis.Get(ctx, key).Int()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				// 键不存在 = 尚无失败记录，属于正常路径
				continue
			}
			global.Log.Error("登录失败计数读取异常，按锁定处理",
				zap.String("key", key), zap.Error(err))
			return true
		}
		if n >= max {
			return true
		}
	}
	return false
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
