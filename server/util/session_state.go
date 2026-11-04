package util

import (
	"errors"
	"fmt"
	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
	"time"
)

type SessionPrincipal struct {
	UserId             string `gorm:"column:id"`
	Username           string
	TenantId           string
	AuthVersion        uint64
	MustChangePassword bool
	IsSuperAdmin       bool
	TenantVersion      uint64 `gorm:"-"`
}

// LoadSessionPrincipal is the common live check for authentication, refresh and requests.
func LoadSessionPrincipal(db *gorm.DB, userID, tenantID string) (*SessionPrincipal, error) {
	if userID == "" || tenantID == "" {
		return nil, errors.New("invalidTenant")
	}
	trusted := TenantBypass(db)
	var user SessionPrincipal
	if tenantID == global.PlatformTenantID {
		if err := trusted.Table("sys_user").Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", userID, tenantID).First(&user).Error; err != nil {
			return nil, principalReadError(err, "userNotFound")
		}
	} else {
		var row struct {
			SessionPrincipal
			TenantPresent     string
			TenantStatus      bool
			TenantExpire      *time.Time
			TenantAuthVersion uint64
			MemberPresent     string
		}
		err := trusted.Table("sys_user").Select("sys_user.id, sys_user.username, sys_user.tenant_id, sys_user.auth_version, sys_user.must_change_password, sys_user.is_super_admin, tenant.tenant_id AS tenant_present, tenant.status AS tenant_status, tenant.expire_time AS tenant_expire, tenant.auth_version AS tenant_auth_version, membership.user_id AS member_present").
			Joins("LEFT JOIN sys_tenant tenant ON tenant.tenant_id = sys_user.tenant_id AND tenant.deleted_at IS NULL").
			Joins("LEFT JOIN sys_user_tenant membership ON membership.user_id = sys_user.id AND membership.tenant_id = sys_user.tenant_id AND membership.status = TRUE AND membership.deleted_at IS NULL").
			Where("sys_user.id = ? AND sys_user.tenant_id = ? AND sys_user.deleted_at IS NULL", userID, tenantID).Take(&row).Error
		if err != nil {
			return nil, principalReadError(err, "userNotFound")
		}
		if row.TenantPresent == "" {
			return nil, errors.New("tenantNotFound")
		}
		if !row.TenantStatus || IsTenantExpired(row.TenantExpire) {
			return nil, errors.New("tenantDisabled")
		}
		if row.MemberPresent == "" {
			return nil, errors.New("userNotFound")
		}
		user = row.SessionPrincipal
		user.TenantVersion = row.TenantAuthVersion
	}
	user.IsSuperAdmin = user.IsSuperAdmin && tenantID == global.PlatformTenantID
	return &user, nil
}

func ValidateSessionVersion(db *gorm.DB, userID, tenantID string, userVersion, tenantVersion uint64) (*SessionPrincipal, error) {
	user, err := LoadSessionPrincipal(db, userID, tenantID)
	if err != nil {
		return nil, err
	}
	if userVersion == 0 || user.AuthVersion != userVersion || user.TenantVersion != tenantVersion {
		return nil, errors.New("refreshTokenNotExist")
	}
	return user, nil
}

var ErrAuthDependencyUnavailable = errors.New("authentication dependency unavailable")

func principalReadError(err error, missing string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(missing)
	}
	return fmt.Errorf("%w: %v", ErrAuthDependencyUnavailable, err)
}
