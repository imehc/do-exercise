package request

type CreateSysUserReq struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=10,startWithLetter,containsLetter"`
	Password string `json:"password" binding:"required,min=6,max=16,complexPassword"`
	Nickname string `json:"nickname"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Avatar   string `json:"avatar"`
	// RoleIds 角色ID，可为空（平台超级管理员创建用户时默认无角色）。
	// 非超管创建时必须至少一个角色，该校验放在 service 层根据操作者身份判断。
	RoleIds []uint `json:"role_ids"`
}

type UpdateSysUserReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	RoleIds  []uint `json:"role_ids"`
}

type UpdateSysUserPasswordReq struct {
	Password string `json:"password" binding:"required,min=6,max=16,complexPassword"`
}

// AssignUserTenantReq 给用户分配（移动）归属租户的请求
type AssignUserTenantReq struct {
	TenantId string `json:"tenant_id" binding:"required"`
}
