package model

type Auth struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	RoleIds      []uint `json:"role_ids"`
	RefreshToken string `json:"refresh_token"`
}
