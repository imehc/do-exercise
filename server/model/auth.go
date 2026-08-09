package model

import "time"

type TokenInfo struct {
	UserId       string    `json:"user_id"`
	Username     string    `json:"username"`
	RoleIds      []uint    `json:"role_ids"`
	RefreshToken string    `json:"refresh_token"`
	Disabled     bool      `json:"disabled"`
	CreatedTime  time.Time `json:"created_time"` // 创建时间
	ExpiredTime  time.Time `json:"expired_time"` // 到期时间
	// MustChangePassword 标记该账号仍需强制修改密码，AuthMiddleware 据此限制可用接口
	MustChangePassword bool `json:"must_change_password"`
}

type RefreshTokenInfo struct {
	UserId      string    `json:"user_id"`
	Username    string    `json:"username"`
	RoleIds     []uint    `json:"role_ids"`
	Disabled    bool      `json:"disabled"`
	CreatedTime time.Time `json:"created_time"`
	ExpiredTime time.Time `json:"expired_time"`
	// MustChangePassword 标记该账号仍需强制修改密码，刷新后的 token 继续携带
	MustChangePassword bool `json:"must_change_password"`
}
