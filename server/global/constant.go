package global

// TenantModeSingle / TenantModeMulti 租户模式枚举
const (
	TenantModeSingle = "single"
	TenantModeMulti  = "multi"
)

// PlatformTenantID 平台层(全局超级管理员)的保留租户 ID。
// sys_tenant 表中不存在该行，仅作为内部标识。
const PlatformTenantID = "platform"

// PlatformOnlyMenuIDs 仅供平台超级管理员使用的菜单（租户管理子树）。
// 创建租户时为租户管理员供应全量菜单时必须排除，业务租户不得获得租户管理权限。
var PlatformOnlyMenuIDs = []uint{10, 181, 182, 183, 184, 185}

// tenantScopedTables 参与行级租户隔离的表。
// 租户插件据此决定是否为当前语句追加 tenant_id 过滤/回填；
// 不在集合内的表（sys_tenant/sys_api/sys_menu/sys_user_role 等）保持全局共享。
var tenantScopedTables = map[string]bool{
	"sys_user":            true,
	"sys_role":            true,
	"sys_job":             true,
	"sys_operation_log":   true,
}

type contextKey string

const (
	// ContextUserIDKey 上下文用户ID
	ContextUserIDKey contextKey = "userId"
	// ContextDBKey 请求绑定的数据库连接
	ContextDBKey contextKey = "gormDB"
	// ContextTenantIDKey 上下文租户ID
	ContextTenantIDKey contextKey = "tenantId"
	// ContextTenantBypassKey 上下文标记：当前请求需要绕过租户隔离。
	// 平台层管理跨租户数据（创建租户/租户管理员等）时由服务显式设置，
	// 租户插件据此跳过 tenant_id 的过滤与回填，改由服务显式控制。
	ContextTenantBypassKey contextKey = "tenantBypass"
)

// IsTenantScopedTable 判断表名是否参与行级租户隔离
func IsTenantScopedTable(table string) bool {
	return tenantScopedTables[table]
}

// ResolveTenantID 依据配置与请求上下文解析当前操作归属的租户ID。
// 返回值 ok=false 表示当前上下文无需租户隔离：
//   - 请求未携带租户（登录等公共端点）
//   - 多租户模式下为平台租户（跨租户管理）
func ResolveTenantID(tid string) (string, bool) {
	if Config.Tenant.IsMulti() {
		if tid == "" || tid == PlatformTenantID {
			return "", false
		}
		return tid, true
	}
	return Config.Tenant.DefaultTenantId, true
}