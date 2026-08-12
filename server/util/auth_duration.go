package util

import "time"

// tokenDurations 启动时一次性解析并校验的 token 时长配置。
// 此前在每次请求里解析字符串，一旦拼写错误，token 生成会拿到 0 值，
// 经 Redis Set(..., 0) 存成永不过期的凭据，而 AuthMiddleware 又默认只信 Redis TTL。
var (
	authAccessExpire  time.Duration
	authRefreshExpire time.Duration
)

// InitAuthDurations 启动时解析并校验 token 时长配置，失败即退出进程。
// 保证后续 AuthMiddleware 与 token 生成拿到的时长始终有效。
func InitAuthDurations(accessExpireTime, refreshExpireTime string) {
	var err error
	authAccessExpire, err = ParseDurationString(accessExpireTime)
	if err != nil {
		Exit("access_expire_time 配置解析失败: ", err)
	}
	authRefreshExpire, err = ParseDurationString(refreshExpireTime)
	if err != nil {
		Exit("refresh_expire_time 配置解析失败: ", err)
	}
}

// AuthAccessExpire 返回已校验的访问令牌时长。
func AuthAccessExpire() time.Duration { return authAccessExpire }

// AuthRefreshExpire 返回已校验的刷新令牌时长。
func AuthRefreshExpire() time.Duration { return authRefreshExpire }