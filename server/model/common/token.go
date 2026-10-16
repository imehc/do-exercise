package common

type Token struct {
	AccessToken       string `json:"access_token"`        // 访问令牌
	ExpireTime        int64  `json:"expire_time"`         // 访问令牌过期时间
	RefreshToken      string `json:"refresh_token"`       // 刷新令牌
	RefreshExpireTime int64  `json:"refresh_expire_time"` // 刷新令牌过期时间
	// MustChangePassword 标记该账号仍需强制修改密码，前端据此强制跳转改密页
	MustChangePassword bool `json:"must_change_password"`
	// TenantId 当前 token 所属租户ID，前端据此定位当前租户（名称经 my_tenants 查询）
	TenantId string `json:"tenant_id"`
	// IsSuperAdmin 平台超级管理员标识。仅平台超管登录时返回 true（omitempty 使
	// 普通账号响应中不含该字段），前端据此区分平台域与业务租户。
	IsSuperAdmin bool `json:"is_super_admin,omitempty"`
}
