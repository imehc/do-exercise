package response

import "github.com/imehc/do-exercise/server/model/common"

// LoginResult 登录/选择租户/切换租户的响应体。
// 多租户模式下账号归属多个启用租户时，requires_tenant_selection=true 且不携带 token，
// 前端展示租户选择后调用 select_tenant 完成登录；否则直接携带 token。
type LoginResult struct {
	*common.Token
	RequiresTenantSelection bool           `json:"requires_tenant_selection"`
	LoginSessionId          string         `json:"login_session_id"`
	AvailableTenants        []TenantOption `json:"available_tenants"`
}
