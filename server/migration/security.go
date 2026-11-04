package migration

import (
	"context"
	"errors"
	"os"
	"unicode"

	"github.com/imehc/do-exercise/server/util"
	"gorm.io/gorm"
)

func hardenTenantSecurity(db *gorm.DB) error {
	if err := execAll(db, []string{
		"DELETE FROM sys_menu_apis WHERE sys_api_id IN (129,130) AND sys_menu_id = 133",
		"INSERT INTO sys_menu_apis (sys_menu_id,sys_api_id) VALUES (183,129),(183,130) ON CONFLICT DO NOTHING",
		"DELETE FROM casbin_rule WHERE ptype = 'p' AND v1 <> 'platform' AND (v2 LIKE '/system/tenants%' OR v2 LIKE '/system/users/%/tenant' OR v2 LIKE '/system/users/%/assignable_tenants' OR v2 = '/system/info')",
		"UPDATE sys_operation_log SET body = '[REDACTED]' WHERE body ~* '(password|refresh_token|access_token|secret)'",
		"UPDATE sys_operation_log SET result = '[REDACTED]', params = '[REDACTED]' WHERE path LIKE '/auth/%' OR path IN ('/system/token/all','/get_presigned_url')",
		"INSERT INTO sys_user_tenant (user_id,tenant_id,status,created_at,updated_at) SELECT id,tenant_id,TRUE,created_at,updated_at FROM sys_user WHERE deleted_at IS NULL AND tenant_id <> '' ON CONFLICT DO NOTHING",
	}); err != nil {
		return err
	}
	var count int64
	if err := db.Table("sys_user").Where("password = ? AND deleted_at IS NULL", defaultAdminHash).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	password := os.Getenv("BOOTSTRAP_ADMIN_PASSWORD")
	if !strongBootstrapPassword(password) {
		return errors.New("set BOOTSTRAP_ADMIN_PASSWORD to a unique 12–72 byte password with upper/lowercase letters, a digit and a symbol before migrating the public default administrator")
	}
	hash, err := (&util.Hash{Value: password}).Hash()
	if err != nil {
		return err
	}
	return db.Table("sys_user").Where("password = ?", defaultAdminHash).Updates(map[string]any{"password": hash, "must_change_password": true, "auth_version": gorm.Expr("auth_version + 1")}).Error
}

func strongBootstrapPassword(password string) bool {
	if len(password) < 12 || len(password) > 72 {
		return false
	}
	var upper, lower, digit, symbol bool
	for _, c := range password {
		upper = upper || unicode.IsUpper(c)
		lower = lower || unicode.IsLower(c)
		digit = digit || unicode.IsDigit(c)
		symbol = symbol || (!unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsSpace(c))
	}
	return upper && lower && digit && symbol
}

// WithSchemaLock serializes AutoMigrate with the versioned migration runner.
func WithSchemaLock(ctx context.Context, db *gorm.DB, fn func() error) error {
	unlock, err := lock(ctx, db)
	if err != nil {
		return err
	}
	defer unlock()
	return fn()
}

// VerifySecurityBaseline prevents serving traffic with an unmigrated or public
// default administrator database, including direct invocation of the server command.
func VerifySecurityBaseline(db *gorm.DB) error {
	var count int64
	if err := db.Table("schema_migrations").Where("version = ?", 11).Count(&count).Error; err != nil || count != 1 {
		return errors.New("run migrate before starting the server (migration 11 is required)")
	}
	if err := util.TenantBypass(db).Table("sys_user").Where("password = ? AND deleted_at IS NULL", defaultAdminHash).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return errors.New("public default administrator password is still present")
	}
	return nil
}
