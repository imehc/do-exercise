package util

import (
	casbinutil "github.com/casbin/casbin/v3/util"
	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
)

// Authorize reads the committed role/menu/API relationships, avoiding a second
// permission store and stale per-process policies. Casbin supplies path matching.
func Authorize(db *gorm.DB, userID, tenantID, path, method string) (bool, error) {
	if userID == "" || tenantID == "" {
		return false, nil
	}
	var apis []struct{ Path string }
	query := db.Table("sys_role").Distinct("sys_api.path").
		Joins("JOIN sys_user_role ON sys_user_role.sys_role_id = sys_role.id").
		Joins("JOIN sys_role_menu ON sys_role_menu.sys_role_id = sys_role.id").
		Joins("JOIN sys_menu ON sys_menu.id = sys_role_menu.sys_menu_id").
		Joins("JOIN sys_menu_apis ON sys_menu_apis.sys_menu_id = sys_menu.id").
		Joins("JOIN sys_api ON sys_api.id = sys_menu_apis.sys_api_id").
		Where("sys_user_role.sys_user_user_id = ? AND sys_role.tenant_id = ? AND sys_role.deleted_at IS NULL", userID, tenantID).
		Where("sys_menu.deleted_at IS NULL AND sys_api.deleted_at IS NULL AND sys_api.disabled = FALSE AND sys_api.method = ?", method)
	if tenantID != global.PlatformTenantID {
		query = query.Where("sys_menu.scope <> ?", global.MenuScopePlatform)
	}
	if err := query.Find(&apis).Error; err != nil {
		return false, err
	}
	for _, api := range apis {
		if casbinutil.KeyMatch2(path, api.Path) {
			return true, nil
		}
	}
	return false, nil
}
