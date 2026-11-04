package system

import (
	"errors"
	"github.com/imehc/do-exercise/server/global"
	sysmodel "github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Caller must use a transaction: the tenant row lock serializes quota checks and writes.
func lockTenant(tx *gorm.DB, tid string) (*sysmodel.SysTenant, error) {
	if tid == global.PlatformTenantID {
		return nil, tx.Exec("SELECT pg_advisory_xact_lock(?)", int64(0x706c61745f61646d)).Error
	}
	if tid == "" {
		return nil, errors.New("invalidTenant")
	}
	var tenant sysmodel.SysTenant
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("tenant_id = ?", tid).First(&tenant).Error; err != nil {
		return nil, err
	}
	if !tenant.Status || util.IsTenantExpired(tenant.ExpireTime) {
		return nil, errors.New("tenantDisabled")
	}
	return &tenant, nil
}
func enforceUserQuota(tx *gorm.DB, tid string, extra int) error {
	tenant, err := lockTenant(tx, tid)
	if err != nil {
		return err
	}
	if tenant == nil || tenant.MaxUsers == 0 || extra <= 0 {
		return nil
	}
	var count int64
	if err := tx.Model(&sysmodel.SysUser{}).Where("tenant_id = ?", tid).Count(&count).Error; err != nil {
		return err
	}
	if count+int64(extra) > int64(tenant.MaxUsers) {
		return errors.New("tenantUserQuotaReached")
	}
	return nil
}
func enforceJobQuota(tx *gorm.DB, tid string, extra int) error {
	tenant, err := lockTenant(tx, tid)
	if err != nil {
		return err
	}
	if tenant == nil || tenant.MaxTasks == 0 || extra <= 0 {
		return nil
	}
	var count int64
	if err := tx.Model(&sysmodel.SysJob{}).Where("tenant_id = ?", tid).Count(&count).Error; err != nil {
		return err
	}
	if count+int64(extra) > int64(tenant.MaxTasks) {
		return errors.New("tenantJobQuotaReached")
	}
	return nil
}

func preserveAdministrator(tx *gorm.DB, user *sysmodel.SysUser, roleIDs []uint, deleting bool) ([]uint, error) {
	if user.IsSuperAdmin && user.TenantId == global.PlatformTenantID && deleting {
		var count int64
		if err := tx.Model(&sysmodel.SysUser{}).Where("tenant_id = ? AND is_super_admin = TRUE AND id <> ?", user.TenantId, user.UserId).Count(&count).Error; err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, errors.New("lastAdministratorRequired")
		}
	}
	var roles []sysmodel.SysRole
	if err := tx.Model(user).Association("Roles").Find(&roles); err != nil {
		return nil, err
	}
	for _, role := range roles {
		if role.Code != TenantAdminRoleCode {
			continue
		}
		retained := false
		for _, id := range roleIDs {
			if id == role.Id {
				retained = true
			}
		}
		if retained {
			return roleIDs, nil
		}
		if !deleting && !isSuperAdmin(tx) {
			return append(roleIDs, role.Id), nil
		}
		var count int64
		if err := tx.Model(&sysmodel.SysUser{}).Joins("JOIN sys_user_role ON sys_user_role.sys_user_user_id = sys_user.id").Where("sys_user_role.sys_role_id = ? AND sys_user.id <> ?", role.Id, user.UserId).Count(&count).Error; err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, errors.New("lastAdministratorRequired")
		}
	}
	return roleIDs, nil
}
