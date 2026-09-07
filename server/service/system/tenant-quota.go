package system

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	sysmodel "github.com/imehc/do-exercise/server/model/system"
	"gorm.io/gorm"
)

// enforceUserQuota 在给租户新增用户前校验用户数上限。
// tenantId 为空或平台保留租户时不限；MaxUsers<=0 表示不设限。
// 配额读取/统计只是尽力而为的守卫：读取或统计出错时按“不限”放行，
// 不让一次临时的配额读取失败阻断本可成功的用户创建。
func enforceUserQuota(db *gorm.DB, tenantId string, extra int) error {
	if tenantId == "" || tenantId == global.PlatformTenantID || extra <= 0 {
		return nil
	}
	var tenant sysmodel.SysTenant
	if err := db.Where("tenant_id = ?", tenantId).First(&tenant).Error; err != nil {
		return nil
	}
	if tenant.MaxUsers <= 0 {
		return nil
	}
	var count int64
	if err := db.Model(&sysmodel.SysUser{}).
		Where("tenant_id = ?", tenantId).
		Count(&count).Error; err != nil {
		return nil
	}
	if count+int64(extra) > int64(tenant.MaxUsers) {
		return errors.New("tenantUserQuotaReached")
	}
	return nil
}

// enforceJobQuota 在给租户新增定时任务前校验任务数上限。口径同 enforceUserQuota。
func enforceJobQuota(db *gorm.DB, tenantId string, extra int) error {
	if tenantId == "" || tenantId == global.PlatformTenantID || extra <= 0 {
		return nil
	}
	var tenant sysmodel.SysTenant
	if err := db.Where("tenant_id = ?", tenantId).First(&tenant).Error; err != nil {
		return nil
	}
	if tenant.MaxTasks <= 0 {
		return nil
	}
	var count int64
	if err := db.Model(&sysmodel.SysJob{}).
		Where("tenant_id = ?", tenantId).
		Count(&count).Error; err != nil {
		return nil
	}
	if count+int64(extra) > int64(tenant.MaxTasks) {
		return errors.New("tenantJobQuotaReached")
	}
	return nil
}
