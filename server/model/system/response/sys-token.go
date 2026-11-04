package response

import "time"

type SysTokenLogRsp struct {
	UserId string `json:"user_id"`
	// TenantId 会话所属租户，供平台超级管理员在全量视图里区分来源
	TenantId            string    `json:"tenant_id"`
	Username            string    `json:"username"`
	SessionId           string    `json:"session_id"`
	Disabled            bool      `json:"disabled"`
	AccessTokenCreated  time.Time `json:"access_token_created"`
	AccessTokenExpired  time.Time `json:"access_token_expired"`
	RefreshTokenCreated time.Time `json:"refresh_token_created"`
	RefreshTokenExpired time.Time `json:"refresh_token_expired"`
}
