package global

import "slices"

// PlatformTenantID 平台层(全局超级管理员)的保留租户 ID。
// sys_tenant 表中不存在该行，仅作为内部标识。
const PlatformTenantID = "platform"

// MenuScope 定义共享菜单的可见范围。
const (
	MenuScopePlatform = "platform"
	MenuScopeTenant   = "tenant"
	MenuScopeBoth     = "both"
)

// MenuPermissionActions 权限标识允许的动作，按界面展示顺序排列。
// 这里是唯一来源：服务端校验 sys_menu.permission 用它，前端的权限动作下拉框
// 也由 GET /auth/tenants 下发同一份，避免两侧各写一套枚举后逐渐漂移。
var MenuPermissionActions = []string{
	"query", "info", "create", "update", "delete",
	"start", "stop", "execute", "reset",
}

// IsMenuPermissionAction 判断动作是否在允许集合内
func IsMenuPermissionAction(action string) bool {
	return slices.Contains(MenuPermissionActions, action)
}

// LegacyPlatformOnlyMenuIDs 旧种子里的平台专属菜单 ID。
//
// 运行期的判定依据已经全部改成 sys_menu.scope（platform/tenant/both），
// 这份名单只剩一个用途：升级存量库时把 scope/is_system 回填到这些行上
// （init.sql 与迁移版本 7 menu_metadata_backfill，见 server/migration）。新增平台菜单只需把 scope 设成 platform，
// 不要再往这里加 ID——自增 ID 不是稳定契约，这也是 P1-3 的整改目标。
//
// 收录标准是「该菜单的数据无法按租户隔离」：
//   - 10/181-185 租户管理：管的就是租户本身，天然跨租户。
//   - 8/161 系统信息：读取的是宿主机 CPU/内存/磁盘/进程，属于平台运维数据，
//     不存在「本租户的 CPU」，暴露给业务租户既无意义又泄露部署拓扑。
//
// 其余菜单（用户/角色/操作记录/令牌/定时任务）的数据都能按租户裁剪，
// 因此保持对租户开放，由行级隔离或服务层范围裁剪保证互不可见。
var LegacyPlatformOnlyMenuIDs = []uint{8, 10, 161, 181, 182, 183, 184, 185}

// tenantScopedTables 参与行级租户隔离的表。
// 租户插件据此决定是否为当前语句追加 tenant_id 过滤/回填；
// 不在集合内的表（sys_tenant/sys_api/sys_menu/sys_user_role 等）保持全局共享。
var tenantScopedTables = map[string]bool{
	"sys_user":          true,
	"sys_role":          true,
	"sys_job":           true,
	"sys_operation_log": true,
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
	// ContextIsSuperAdminKey 上下文标记：当前操作者是平台超级管理员。
	// 由 ContextMiddleware 从已通过校验的会话（AuthMiddleware 写入的 gin 键）透传，
	// 服务层据此判断可见范围是否放开到全部租户，避免每次都回查 sys_user。
	ContextIsSuperAdminKey contextKey = "isSuperAdmin"
)

// IsTenantScopedTable 判断表名是否参与行级租户隔离
func IsTenantScopedTable(table string) bool {
	return tenantScopedTables[table]
}

// ResolveTenantID 依据请求上下文解析当前操作归属的租户ID。
// 系统只有多租户一种形态，返回值 ok=false 表示当前上下文无需租户隔离：
//   - 请求未携带租户（登录等公共端点、迁移、定时任务的平台维护旁路）
//   - 平台租户（跨租户管理）
func ResolveTenantID(tid string) (string, bool) {
	if tid == "" || tid == PlatformTenantID {
		return "", false
	}
	return tid, true
}
