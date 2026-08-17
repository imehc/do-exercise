package request

// CreateSysTenantReq 创建租户请求。
// 管理员账号二选一：
//   - admin_mode = new（默认）：新建用户，需 admin_username + admin_password
//   - admin_mode = existing：选择现有用户（排除平台超级管理员），需 admin_user_id
type CreateSysTenantReq struct {
	Name   string `json:"name" binding:"required,min=2,max=64"`
	Code   string `json:"code" binding:"required,alphanum,min=2,max=32"`
	Remark string `json:"remark,omitempty" binding:"omitempty,max=255"`
	// AdminMode 管理员账号模式：new / existing
	AdminMode string `json:"admin_mode" binding:"required,oneof=new existing"`
	// AdminUserId 选择现有用户时的用户 ID（admin_mode=existing）
	AdminUserId string `json:"admin_user_id,omitempty" binding:"omitempty"`
	// AdminUsername / AdminPassword 新建用户时的账号信息（admin_mode=new）
	AdminUsername string `json:"admin_username,omitempty" binding:"omitempty,alphanum,min=2,max=10,startWithLetter,containsLetter"`
	AdminPassword string `json:"admin_password,omitempty" binding:"omitempty,min=6,max=16,complexPassword"`
}

// AdminMode 常量
const (
	AdminModeNew      = "new"
	AdminModeExisting = "existing"
)

type UpdateSysTenantReq struct {
	Name   string `json:"name" binding:"required,min=2,max=64"`
	Status *bool  `json:"status"`
	Remark string `json:"remark,omitempty" binding:"omitempty,max=255"`
}

type QuerySysTenantReq struct {
	Name  string `json:"name" form:"name"`
	Code  string `json:"code" form:"code"`
	Page  int    `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

// AssignTenantUsersReq 分配现有用户到租户的请求。
// user_ids 为用户 ID（sys_user.id）列表，非空且去重后逐个复制到目标租户下。
type AssignTenantUsersReq struct {
	UserIds []string `json:"user_ids" binding:"required,min=1,dive,required"`
}