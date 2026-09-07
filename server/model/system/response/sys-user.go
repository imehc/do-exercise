package response

import (
	"time"
)

type SysUserResp struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	// TenantId 用户所属租户，用于管理员确认「这个人在哪个租户」
	TenantId string `json:"tenant_id"`
	// TenantName 租户名称，由服务层批量回填（sys_tenant 不参与行级隔离，可直接查）；
	// 租户已被删除或数据缺失时为空串，前端回退展示 TenantId。
	TenantName string        `json:"tenant_name"`
	CreatedAt  time.Time     `json:"created_at,omitzero"`
	CreatedBy  int64         `json:"created_by,omitzero"`
	UpdatedAt  time.Time     `json:"updated_at,omitzero"`
	UpdatedBy  int64         `json:"updated_by,omitzero"`
	Roles      []SysRoleResp `json:"roles"`
}

// AssignableTenant 用户列表中可分配的候选租户
type AssignableTenant struct {
	TenantId string `json:"tenant_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
}
