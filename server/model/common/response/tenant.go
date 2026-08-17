package response

// TenantOption 登录页租户选择项
type TenantOption struct {
	TenantId string `json:"tenant_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
}

// AvailableTenants 登录页可用的租户列表（含当前模式，前端据此决定是否渲染选择器）
type AvailableTenants struct {
	Mode    string         `json:"mode"`
	Tenants []TenantOption `json:"tenants"`
}