package storage

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const SignatureLifetime = 5 * time.Minute
const UploadGrace = 5 * time.Minute

func tenantLimit(tx *gorm.DB, tenantID string) (int64, error) {
	if tenantID == global.PlatformTenantID {
		return 0, nil
	}
	var tenant system.SysTenant
	if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).First(&tenant, "tenant_id = ?", tenantID).Error; err != nil {
		return 0, err
	}
	if !tenant.Status || util.IsTenantExpired(tenant.ExpireTime) {
		return 0, errors.New("tenantDisabled")
	}
	return tenant.MaxStorageBytes, nil
}

// Locking the counter makes reservations atomic across all application replicas.
func RegisterExisting(tx *gorm.DB, file *system.SysFile) error { return charge(tx, file, 0) }

func Reserve(tx *gorm.DB, file *system.SysFile) error {
	limit, err := tenantLimit(tx, file.TenantId)
	if err != nil {
		return err
	}
	return charge(tx, file, limit)
}
func charge(tx *gorm.DB, file *system.SysFile, limit int64) error {
	if file.MaxSize <= 0 {
		return errors.New("invalidFileSize")
	}
	usage := system.SysStorageUsage{TenantId: file.TenantId}
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&usage).Error; err != nil {
		return err
	}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&usage, "tenant_id = ?", file.TenantId).Error; err != nil {
		return err
	}
	if usage.ChargedBytes > math.MaxInt64-file.MaxSize || limit > 0 && (file.MaxSize > limit || usage.ChargedBytes > limit-file.MaxSize) {
		return errors.New("tenantStorageQuotaReached")
	}
	if err := tx.Model(&usage).UpdateColumn("charged_bytes", gorm.Expr("charged_bytes + ?", file.MaxSize)).Error; err != nil {
		return err
	}
	return tx.Create(file).Error
}
func avatarReferences(tx *gorm.DB, key string) (int64, error) {
	var count int64
	// Copied accounts may still reference an image from their source tenant.
	err := util.TenantBypass(tx).Model(&system.SysUser{}).Where("avatar IN ?", []string{"/" + global.Config.Oss.BucketName + "/" + key, "/oss/" + global.Config.Oss.BucketName + "/" + key}).Count(&count).Error
	return count, err
}
func Complete(db *gorm.DB, userID, key string) (*system.SysFile, error) {
	var file system.SysFile
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&file, "object_key = ? AND user_id = ?", key, userID).Error; err != nil {
			return errors.New("fileNotFound")
		}
		if file.State == "ready" {
			return nil
		}
		if file.State != "pending" || file.UploadExpiresAt == nil || time.Now().After(file.UploadExpiresAt.Add(UploadGrace)) {
			return errors.New("uploadExpired")
		}
		info, err := global.Oss.StatObject(tx.Statement.Context, global.Config.Oss.BucketName, key, minio.StatObjectOptions{})
		if err != nil || info.Size != file.MaxSize || info.ContentType != file.ContentType {
			return errors.New("invalidUploadedFile")
		}
		now := time.Now().UTC()
		file.Size = info.Size
		file.State = "ready"
		file.CompletedAt = &now
		return tx.Model(&file).UpdateColumns(map[string]any{"size": file.Size, "state": file.State, "completed_at": now}).Error
	})
	return &file, err
}
func Delete(db *gorm.DB, userID, key string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var file system.SysFile
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&file, "object_key = ? AND user_id = ?", key, userID).Error; err != nil {
			return errors.New("fileNotFound")
		}
		if file.State == "deleted" || file.State == "deleting" {
			return nil
		}
		if file.Public {
			count, err := avatarReferences(tx, key)
			if err != nil {
				return err
			}
			if count > 0 {
				return errors.New("fileInUse")
			}
		}
		after := time.Now().Add(SignatureLifetime + UploadGrace)
		if file.UploadExpiresAt != nil && after.Before(file.UploadExpiresAt.Add(UploadGrace)) {
			after = file.UploadExpiresAt.Add(UploadGrace)
		}
		return tx.Model(&file).UpdateColumns(map[string]any{"state": "deleting", "delete_after": after}).Error
	})
}
func Usage(db *gorm.DB, tenantID string) (int64, int64, error) {
	limit := int64(0)
	if tenantID != global.PlatformTenantID {
		var tenant system.SysTenant
		if err := db.First(&tenant, "tenant_id = ?", tenantID).Error; err != nil {
			return 0, 0, err
		}
		limit = tenant.MaxStorageBytes
	}
	var usage system.SysStorageUsage
	err := db.First(&usage, "tenant_id = ?", tenantID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return usage.ChargedBytes, limit, err
}

// Collect removes only expired reservations, explicitly deleted files, unused
// avatars and files belonging to tenants deleted beyond their retention period.
// Objects are removed before their quota is released. Failed deletes stay charged.
func Collect(ctx context.Context, db *gorm.DB, now time.Time, limit int) (int, error) {
	if limit < 1 || limit > 1000 {
		limit = 100
	}
	db = util.TenantBypass(db.WithContext(ctx))
	avatarHours := global.Config.Oss.AvatarRetentionHours
	if avatarHours <= 0 {
		avatarHours = 24
	}
	tenantDays := global.Config.Oss.DeletedTenantRetentionDays
	if tenantDays <= 0 {
		tenantDays = 30
	}
	var files []system.SysFile
	if err := db.Where("state <> 'deleted' AND state <> 'legacy'").Where(`(state = 'pending' AND upload_expires_at < ?) OR (state = 'deleting' AND delete_after <= ?) OR (state = 'ready' AND public = TRUE AND completed_at < ? AND NOT EXISTS (SELECT 1 FROM sys_user WHERE deleted_at IS NULL AND avatar IN ('/' || ? || '/' || sys_file.object_key, '/oss/' || ? || '/' || sys_file.object_key))) OR tenant_id IN (SELECT tenant_id FROM sys_tenant WHERE deleted_at < ?)`, now.Add(-UploadGrace), now, now.Add(-time.Duration(avatarHours)*time.Hour), global.Config.Oss.BucketName, global.Config.Oss.BucketName, now.Add(-time.Duration(tenantDays)*24*time.Hour)).Order("created_at,object_key").Limit(limit).Find(&files).Error; err != nil {
		return 0, err
	}
	deleted := 0
	var failures []error
	for _, snapshot := range files {
		err := db.Transaction(func(tx *gorm.DB) error {
			var file system.SysFile
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).First(&file, "object_key = ?", snapshot.ObjectKey).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			} else if err != nil {
				return err
			}
			if file.State == "deleted" || file.State == "legacy" {
				return nil
			}
			var tenant system.SysTenant
			err := tx.Unscoped().First(&tenant, "tenant_id = ?", file.TenantId).Error
			tenantDeleted := err == nil && tenant.DeletedAt.Valid && tenant.DeletedAt.Time.Before(now.Add(-time.Duration(tenantDays)*24*time.Hour))
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			due := file.State == "pending" && file.UploadExpiresAt != nil && file.UploadExpiresAt.Add(UploadGrace).Before(now) || file.State == "deleting" && file.DeleteAfter != nil && !file.DeleteAfter.After(now)
			if !due && !tenantDeleted {
				if file.State != "ready" || !file.Public || file.CompletedAt == nil || !file.CompletedAt.Before(now.Add(-time.Duration(avatarHours)*time.Hour)) {
					return nil
				}
			}
			if file.UploadExpiresAt != nil && now.Before(file.UploadExpiresAt.Add(UploadGrace)) {
				return nil
			}
			if file.Public {
				count, err := avatarReferences(tx, file.ObjectKey)
				if err != nil {
					return err
				}
				if count > 0 {
					return nil
				}
			}
			if err := global.Oss.RemoveObject(ctx, global.Config.Oss.BucketName, file.ObjectKey, minio.RemoveObjectOptions{}); err != nil {
				return fmt.Errorf("delete %s: %w", file.ObjectKey, err)
			}
			result := tx.Model(&system.SysStorageUsage{}).Where("tenant_id = ? AND charged_bytes >= ?", file.TenantId, file.MaxSize).UpdateColumn("charged_bytes", gorm.Expr("charged_bytes - ?", file.MaxSize))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return errors.New("storage counter inconsistent")
			}
			if err := tx.Model(&file).UpdateColumns(map[string]any{"state": "deleted", "deleted_at": now}).Error; err != nil {
				return err
			}
			deleted++
			return nil
		})
		if err != nil {
			failures = append(failures, err)
		}
	}
	return deleted, errors.Join(failures...)
}

// LockAvatar keeps a newly attached local avatar alive through the user update.
func LockAvatar(tx *gorm.DB, avatar, tenantID string) error {
	prefix := "/" + global.Config.Oss.BucketName + "/public/avatars/"
	normalized := strings.TrimPrefix(avatar, "/oss")
	if !strings.HasPrefix(normalized, prefix) {
		return nil
	}
	key := strings.TrimPrefix(normalized, "/"+global.Config.Oss.BucketName+"/")
	var file system.SysFile
	if err := util.TenantBypass(tx).Clauses(clause.Locking{Strength: "SHARE"}).First(&file, "object_key = ? AND tenant_id = ? AND public = TRUE AND state = 'ready'", key, tenantID).Error; err != nil {
		return errors.New("fileNotFound")
	}
	return nil
}

// Canonical local paths keep avatar references recognizable by the collector.
func NormalizeAvatar(avatar string) string {
	parsed, err := url.Parse(avatar)
	if err != nil {
		return avatar
	}
	if parsed.Host != "" {
		public, _ := url.Parse(global.Config.Oss.PresignedHost)
		local := fmt.Sprintf("%s:%d", global.Config.Oss.Host, global.Config.Oss.Port)
		if parsed.Host != local && (public == nil || public.Host == "" || parsed.Host != public.Host) {
			return avatar
		}
	}
	path := strings.TrimPrefix(parsed.Path, "/oss")
	if strings.HasPrefix(path, "/"+global.Config.Oss.BucketName+"/public/avatars/") {
		return path
	}
	return avatar
}
