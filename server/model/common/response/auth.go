package response

import "github.com/imehc/do-exercise/server/model/common"

// LoginResult 登录/选择租户/切换租户的响应体。
// 账号归属多个启用租户时，requires_tenant_selection=true 且不携带 token，
// 前端展示租户选择后调用 select_tenant 完成登录；否则直接携带 token。
type LoginResult struct {
	*common.Token
	RequiresTenantSelection bool           `json:"requires_tenant_selection"`
	LoginSessionId          string         `json:"login_session_id"`
	AvailableTenants        []TenantOption `json:"available_tenants"`
}

// ResetPasswordResult 找回密码的响应体。
// 同一邮箱归属多个租户时 requires_tenant_selection=true 且**未做任何修改**，
// 前端展示租户选择后带 tenant_id 重试（该次验证码仍然有效，不必重新收信）。
type ResetPasswordResult struct {
	RequiresTenantSelection bool           `json:"requires_tenant_selection"`
	AvailableTenants        []TenantOption `json:"available_tenants"`
}
