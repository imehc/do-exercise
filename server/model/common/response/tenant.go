package response

// TenantOption 登录页租户选择项
type TenantOption struct {
	TenantId string `json:"tenant_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
}

// AvailableTenants 登录页引导数据。
// 系统只有多租户一种形态，不再下发部署模式；tenants 也始终为空（禁止匿名枚举租户），
// 认证后的候选租户由登录结果返回。
type AvailableTenants struct {
	Tenants []TenantOption `json:"tenants"`
	// PermissionActions 权限标识允许的动作枚举，由服务端集中下发。
	// 前端菜单管理的「权限动作」下拉框直接消费这份列表，不再各写一套常量。
	// 这是一份静态词表（query/create/...），不含任何部署或租户信息，
	// 因此放在这个匿名可访问的引导端点上是安全的。
	PermissionActions []string `json:"permission_actions"`
}
