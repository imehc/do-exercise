package system

import (
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// 本文件集中放置「租户可见范围」的判定与查询裁剪，避免各服务各写一套。
//
// 行级隔离（sys_user / sys_role / sys_job / sys_operation_log）由 GORM 租户插件
// 自动完成，服务层无需干预。但有三类范围问题插件管不到，必须显式处理：
//
//  1. 全局定义表（sys_menu / sys_api）没有 tenant_id，插件不会过滤，
//     必须按「平台专属」名单裁剪可见集合，并禁止租户改写这些共享定义。
//  2. Redis 里的会话（令牌管理）不经过 GORM，需要按 TokenInfo.TenantId 自行过滤。
//  3. 租户管理员角色（tenant_admin）虽在本租户内，但它是平台供应的内建角色，
//     不能被租户自己改写或删除。

// isSuperAdmin 判断当前操作者是否平台超级管理员。
func isSuperAdmin(db *gorm.DB) bool {
	return model.CurrentIsSuperAdmin(db)
}

// tenantRestricted 判断当前请求是否受租户可见范围限制。
// 仅平台超级管理员不受限——按设计他需要跨租户视图；其余调用者一律受限。
func tenantRestricted(db *gorm.DB) bool {
	return !isSuperAdmin(db)
}

// scopeTenantVisibleMenus 为受限调用者裁掉平台专属菜单（租户管理子树）。
// sys_menu 是跨租户共享的定义表，没有 tenant_id，因此「当前租户下的菜单」
// 落地为「该租户可持有的菜单」——按 scope=platform 排除。
//
// 这里不采用「该租户各角色已授予菜单的并集」：菜单管理对租户是只读的，
// 一旦某菜单被从角色上摘掉就再也看不到、也无法重新授予，会把租户锁死。
//
// 判定只写 `scope <> platform`，不再带 `scope IS NULL OR` 兜底：那种兜底是
// fail-open 的（NULL 会被当成租户可见），而 scope 列已由迁移版本 8 约束成
// NOT NULL DEFAULT 'both'，NULL 不再是一种可能的状态。
func scopeTenantVisibleMenus(db *gorm.DB) *gorm.DB {
	if !tenantRestricted(db) {
		return db
	}
	return db.Where("scope <> ?", global.MenuScopePlatform)
}

// FilterTenantVisibleMenus 对已加载的菜单集合应用与 scopeTenantVisibleMenus 相同的规则。
// /user/menu 等基于角色关联预加载的路径无法直接套查询 scope，统一走这里避免口径漂移。
func FilterTenantVisibleMenus(db *gorm.DB, menus []system.SysMenu) []system.SysMenu {
	if !tenantRestricted(db) {
		return menus
	}
	return lo.Filter(menus, func(menu system.SysMenu, _ int) bool {
		return !IsPlatformOnlyMenu(menu)
	})
}

// IsPlatformOnlyMenu 判断菜单是否仅平台可用。
// 唯一依据是显式的 scope 列；旧种子的固定 ID 名单已由 init.sql 与迁移版本 7
// （server/migration，menu_metadata_backfill）回填成 scope=platform，
// 运行期不再按 ID 猜测（见 global.LegacyPlatformOnlyMenuIDs）。
func IsPlatformOnlyMenu(menu system.SysMenu) bool {
	return menu.Scope == global.MenuScopePlatform
}

// scopeTenantVisibleApis 为受限调用者裁掉平台专属 API。
// 判定口径：只被平台专属菜单绑定、且不被任何租户可见菜单绑定的 API 不可见。
// 同时挂在两侧的 API 仍可见；未绑定任何菜单的公共端点（登录、验证码等）也保持可见，
// 否则租户管理员在 API 列表里会看不到自己每天都在用的接口。
func scopeTenantVisibleApis(db *gorm.DB) *gorm.DB {
	if !tenantRestricted(db) {
		return db
	}
	// 平台专属菜单，按 scope 判定（与 IsPlatformOnlyMenu 同一口径）
	platformMenus := db.Session(&gorm.Session{NewDB: true}).
		Table("sys_menu").
		Select("id").
		Where("scope = ?", global.MenuScopePlatform)
	// 租户可见菜单绑定的 API
	tenantApis := db.Session(&gorm.Session{NewDB: true}).
		Table(menuApiJoinTable).
		Select("sys_api_id").
		Where("sys_menu_id NOT IN (?)", platformMenus)
	// 仅被平台专属菜单绑定的 API
	platformOnlyApis := db.Session(&gorm.Session{NewDB: true}).
		Table(menuApiJoinTable).
		Select("sys_api_id").
		Where("sys_menu_id IN (?)", platformMenus).
		Where("sys_api_id NOT IN (?)", tenantApis)
	return db.Where("id NOT IN (?)", platformOnlyApis)
}

// menuApiJoinTable 菜单与 API 的多对多关联表，与 system.SysMenu.Apis 的
// `many2many:sys_menu_apis` 标签保持一致。
const menuApiJoinTable = "sys_menu_apis"

// isTenantAdminRole 判断角色是否为平台供应的租户管理员内建角色。
func isTenantAdminRole(role *system.SysRole) bool {
	return role != nil && role.Code == TenantAdminRoleCode
}
