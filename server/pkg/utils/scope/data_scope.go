package scope

import (
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"gorm.io/gorm"
)

// GetDataScope 获取数据范围
func GetDataScope(db *gorm.DB, user *common.ScopeData, tableName string) *gorm.DB {
	if user == nil || user.RoleId == 0 {
		return db
	}

	// 获取用户角色
	var role system.Role
	if err := db.First(&role, user.RoleId).Error; err != nil {
		return db
	}

	// admin角色不受数据权限限制
	if role.IsAdmin {
		return db
	}

	// TODO: 最后实现该功能
	return db

	// // 根据角色的数据范围构建查询条件
	// switch role.DataScope {
	// case 1: // 全部数据权限
	// 	return db
	// case 2: // 自定数据权限
	// 	return db.Where(tableName+".create_by in (select sys_user.id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", user.RoleId)
	// case 3: // 本部门数据权限
	// 	return db.Where(tableName+".create_by in (SELECT id from sys_user where dept_id = ? )", user.DeptId)
	// case 4: // 本部门及以下数据权限
	// 	return db.Where(tableName+".create_by in (SELECT id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where path like ? ))", "%"+fmt.Sprintf("/%d", user.DeptId)+"%")
	// // case 5: // 仅本人数据权限
	// default:
	// 	return db.Where(tableName+".create_by = ?", user.UserId)
	// }
}
