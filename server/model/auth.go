package model

type Auth struct {
	UserID  int64  `json:"user_id"`
	RoleIds []uint `json:"role_ids"`
}
