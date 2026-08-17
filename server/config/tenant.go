package config

// Tenant 多租户配置
type Tenant struct {
	// Mode 租户模式：single(单租户，固定使用 DefaultTenantId) / multi(多租户)
	Mode string `yaml:"mode" mapstructure:"mode"`
	// DefaultTenantId 单租户模式下所有数据的归属租户，也作为多租户登录未指定租户时的兜底
	DefaultTenantId string `yaml:"default_tenant_id" mapstructure:"default_tenant_id"`
}

// IsMulti 是否多租户模式
func (t Tenant) IsMulti() bool {
	return t.Mode == "multi"
}
