package util

import "time"

// IsTenantExpired 判断租户是否已到期（ExpireTime 非空且早于当前时间）。
// 用于「过期视同停用」的判定，供登录、请求态复核与租户更新时统一复用。
func IsTenantExpired(expireTime *time.Time) bool {
	return expireTime != nil && !time.Now().Before(*expireTime)
}
