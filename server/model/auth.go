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
}

type RefreshTokenInfo struct {
	UserId      string    `json:"user_id"`
	Username    string    `json:"username"`
	RoleIds     []uint    `json:"role_ids"`
	Disabled    bool      `json:"disabled"`
	CreatedTime time.Time `json:"created_time"`
	ExpiredTime time.Time `json:"expired_time"`
}
