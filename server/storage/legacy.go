package storage

import (
	"context"
	"errors"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Old signatures can still overwrite their key until expiry. Only after that
// window may the conservative reservation shrink to the actual object size.
func ReconcileLegacy(ctx context.Context, db *gorm.DB, now time.Time, limit int) error {
	db = util.TenantBypass(db.WithContext(ctx))
	if limit < 1 || limit > 1000 {
		limit = 100
	}
	var files []system.SysFile
	if err := db.Where("state = 'legacy' AND created_at < ?", now.Add(-SignatureLifetime-UploadGrace)).Order("created_at,object_key").Limit(limit).Find(&files).Error; err != nil {
		return err
	}
	var failures []error
	for _, snapshot := range files {
		err := db.Transaction(func(tx *gorm.DB) error {
			var file system.SysFile
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).First(&file, "object_key = ? AND state = 'legacy'", snapshot.ObjectKey).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			} else if err != nil {
				return err
			}
			info, err := global.Oss.StatObject(ctx, global.Config.Oss.BucketName, file.ObjectKey, minio.StatObjectOptions{})
			if err != nil {
				code := minio.ToErrorResponse(err).Code
				if code != "NoSuchKey" && code != "NoSuchObject" && code != "NotFound" {
					return err
				}
				expired := now.Add(-UploadGrace - time.Second)
				return tx.Model(&file).UpdateColumns(map[string]any{"state": "pending", "upload_expires_at": expired}).Error
			}
			if info.Size <= 0 || info.Size > file.MaxSize || info.ContentType != file.ContentType {
				return errors.New("legacy object metadata mismatch")
			}
			if err := tx.Model(&system.SysStorageUsage{}).Where("tenant_id = ? AND charged_bytes >= ?", file.TenantId, file.MaxSize-info.Size).UpdateColumn("charged_bytes", gorm.Expr("charged_bytes - ?", file.MaxSize-info.Size)).Error; err != nil {
				return err
			}
			return tx.Model(&file).UpdateColumns(map[string]any{"state": "ready", "size": info.Size, "max_size": info.Size, "completed_at": file.CreatedAt}).Error
		})
		if err != nil {
			failures = append(failures, err)
		}
	}
	return errors.Join(failures...)
}
