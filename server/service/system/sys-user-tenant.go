package system

import (
	sysmodel "github.com/imehc/do-exercise/server/model/system"
	"gorm.io/gorm"
)

// addUserTenantMembership 记录「userId 成为 tenantId 租户的成员」这一事实。
// 与 sys_user 行的创建在同一个事务里调用，保证两个表要么都写入、要么都回滚，
// 不让成员关系表和用户表出现半同步状态。
func addUserTenantMembership(tx *gorm.DB, userId, tenantId string) error {
	if userId == "" || tenantId == "" {
		return nil
	}
	return tx.Create(&sysmodel.SysUserTenant{
		UserId:   userId,
		TenantId: tenantId,
		Status:   true,
	}).Error
}

// removeUserTenantMembership 撤销成员关系（软删除），与 sys_user 行的移出在同一事务里。
// 已移出租户的用户行是软删除的，成员关系随之软删除，保持两边口径一致。
func removeUserTenantMembership(tx *gorm.DB, userId, tenantId string) error {
	if userId == "" || tenantId == "" {
		return nil
	}
	return tx.
		Where("user_id = ? AND tenant_id = ?", userId, tenantId).
		Delete(&sysmodel.SysUserTenant{}).
		Error
}
