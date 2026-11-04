package util

import (
	"context"
	"errors"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
)

// TenantBypass is for trusted internal operations (authentication, migrations and maintenance).
// HTTP platform operations must validate their principal before obtaining this connection.
func TenantBypass(db *gorm.DB) *gorm.DB {
	return db.WithContext(context.WithValue(db.Statement.Context, global.ContextTenantBypassKey, true))
}

func RequirePlatform(db *gorm.DB) error {
	tid, _ := db.Statement.Context.Value(global.ContextTenantIDKey).(string)
	super, _ := db.Statement.Context.Value(global.ContextIsSuperAdminKey).(bool)
	if tid != global.PlatformTenantID || !super {
		return errors.New("platformOnly")
	}
	return nil
}

func RequireTenantUsable(db *gorm.DB, tenantID string) error {
	if tenantID == global.PlatformTenantID {
		return nil
	}
	if tenantID == "" {
		return errors.New("invalidTenant")
	}
	var tenant struct {
		Status     bool
		ExpireTime *time.Time
	}
	if err := TenantBypass(db).Table("sys_tenant").Where("tenant_id = ? AND deleted_at IS NULL", tenantID).First(&tenant).Error; err != nil {
		return errors.New("tenantNotFound")
	}
	if !tenant.Status || IsTenantExpired(tenant.ExpireTime) {
		return errors.New("tenantDisabled")
	}
	return nil
}
