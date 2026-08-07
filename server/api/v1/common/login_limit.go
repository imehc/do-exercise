package common

import (
	"context"
	"time"

	"github.com/imehc/do-exercise/server/global"
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

// loginIsLocked 判断当前 IP/用户名是否已被锁定
func loginIsLocked(ip, username string) bool {
	ctx := context.Background()
	max, _ := loginAttempts()
	for _, key := range loginFailKeys(ip, username) {
		n, err := global.Redis.Get(ctx, key).Int()
		if err == nil && n >= max {
			return true
		}
	}
	return false
}

// registerLoginFailure 记录一次登录失败，首次失败时设定过期时间
func registerLoginFailure(ip, username string) {
	ctx := context.Background()
	_, lock := loginAttempts()
	for _, key := range loginFailKeys(ip, username) {
		n, err := global.Redis.Incr(ctx, key).Result()
		if err != nil {
			continue
		}
		if n == 1 {
			global.Redis.Expire(ctx, key, lock)
		}
	}
}

// clearLoginFailures 登录成功后清除失败计数
func clearLoginFailures(ip, username string) {
	ctx := context.Background()
	for _, key := range loginFailKeys(ip, username) {
		global.Redis.Del(ctx, key)
	}
}
