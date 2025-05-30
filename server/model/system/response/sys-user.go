package response

import (
	"time"
)

type SysUserResp struct {
	Id        string        `json:"id"`
	Username  string        `json:"username"`
	Nickname  string        `json:"nickname"`
	Email     string        `json:"email"`
	Avatar    string        `json:"avatar"`
	CreatedAt time.Time     `json:"created_at,omitzero"`
	CreatedBy int64         `json:"created_by,omitzero"`
	UpdatedAt time.Time     `json:"updated_at,omitzero"`
	UpdatedBy int64         `json:"updated_by,omitzero"`
	Roles     []SysRoleResp `json:"roles"`
}
