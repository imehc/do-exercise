package model

import "time"

type TokenInfo struct {
	UserId       string    `json:"user_id"`
	Username     string    `json:"username"`
	RoleIds      []uint    `json:"role_ids"`
	TenantId     string    `json:"tenant_id"`
	RefreshToken string    `json:"refresh_token"`
	Disabled     bool      `json:"disabled"`
	CreatedTime  time.Time `json:"created_time"` // 创建时间
	ExpiredTime  time.Time `json:"expired_time"` // 到期时间
	// MustChangePassword 标记该账号仍需强制修改密码，AuthMiddleware 据此限制可用接口
	MustChangePassword bool `json:"must_change_password"`
	// IsSuperAdmin 平台超级管理员标识，AuthMiddleware/Casbin 据此放行平台域管理接口
	IsSuperAdmin bool `json:"is_super_admin"`
}

type RefreshTokenInfo struct {
	UserId      string    `json:"user_id"`
	Username    string    `json:"username"`
	RoleIds     []uint    `json:"role_ids"`
	TenantId    string    `json:"tenant_id"`
	Disabled    bool      `json:"disabled"`
	CreatedTime time.Time `json:"created_time"`
	ExpiredTime time.Time `json:"expired_time"`
	// MustChangePassword 标记该账号仍需强制修改密码，刷新后的 token 继续携带
	MustChangePassword bool `json:"must_change_password"`
	// IsSuperAdmin 平台超级管理员标识，刷新后的 token 继续携带
	IsSuperAdmin bool `json:"is_super_admin"`
	// Rotated 标记该 refresh token 已被轮转消费。
	// 轮转后旧值保留（仅标记），再次出现即判定为家族失陷并吊销该用户全部会话。
	Rotated bool `json:"rotated"`
}
