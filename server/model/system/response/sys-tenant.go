package response

import "time"

type SysTenantResp struct {
	TenantId   string     `json:"tenant_id"`
	Name       string     `json:"name"`
	Code       string     `json:"code"`
	Status     bool       `json:"status"`
	ExpireTime *time.Time `json:"expire_time,omitzero"`
	MaxUsers   int        `json:"max_users"`
	MaxTasks   int        `json:"max_tasks"`
	Remark     string     `json:"remark"`
	CreatedAt  time.Time  `json:"created_at,omitzero"`
	CreatedBy  string     `json:"created_by,omitzero"`
	UpdatedAt  time.Time  `json:"updated_at,omitzero"`
	UpdatedBy  string     `json:"updated_by,omitzero"`
}

// AssignableUser 可分配给租户的现有用户（平台超级管理员与已归属目标租户的用户不在此列）。
type AssignableUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email,omitempty"`
	TenantId string `json:"tenant_id"`
}
