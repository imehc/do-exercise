package internal

import (
	"context"
	"crypto/sha256"
	"fmt"
	"mime"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/storage"
	"github.com/imehc/do-exercise/server/util"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// MigrateLegacyAvatars copies the old publicly referenced images once, before
// accepting requests. All other legacy objects remain private. Copies use stable
// keys so a failed migration can safely retry without exposing the old namespace.
func MigrateLegacyAvatars(ctx context.Context, client *minio.Client) error {
	db := util.TenantBypass(global.DB.WithContext(ctx))
	if err := db.Exec("CREATE TABLE IF NOT EXISTS sys_storage_migrations (name TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())").Error; err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", int64(0x6f73735f6d6967)).Error; err != nil {
			return err
		}
		var count int64
		if err := tx.Table("sys_storage_migrations").Where("name = ?", "private_objects_v1").Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return nil
		}
		var users []system.SysUser
		if err := tx.Unscoped().Find(&users).Error; err != nil {
			return err
		}
		bucket := global.Config.Oss.BucketName
		for _, user := range users {
			// External avatars and new uploads are not migration candidates.
			parsed, err := url.Parse(user.Avatar)
			if err != nil {
				continue
			}
			if parsed.Host != "" {
				continue
			}
			legacy := strings.TrimPrefix(parsed.Path, "/oss")
			if !strings.HasPrefix(legacy, "/"+bucket+"/") {
				continue
			}
			key := strings.TrimPrefix(legacy, "/"+bucket+"/")
			if !strings.HasPrefix(key, user.TenantId+"/"+user.UserId+"/") {
				continue
			}
			ext := strings.ToLower(path.Ext(key))
			allowed := map[string]string{".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png", ".gif": "image/gif", ".webp": "image/webp"}
			contentType, ok := allowed[ext]
			if !ok {
				continue
			}
			info, err := client.StatObject(ctx, bucket, key, minio.StatObjectOptions{})
			if err != nil {
				if minio.ToErrorResponse(err).Code == "NoSuchKey" {
					continue
				}
				return err
			}
			if info.Size <= 0 || info.Size > 10<<20 {
				continue
			}
			digest := sha256.Sum256([]byte(key))
			next := fmt.Sprintf("public/avatars/%s/%s/legacy-%x%s", user.TenantId, user.UserId, digest[:16], ext)
			if _, err := client.CopyObject(ctx, minio.CopyDestOptions{Bucket: bucket, Object: next, ReplaceMetadata: true, UserMetadata: map[string]string{"Content-Type": contentType}}, minio.CopySrcOptions{Bucket: bucket, Object: key}); err != nil {
				return err
			}
			metadata := system.SysFile{ObjectKey: next, TenantId: user.TenantId, UserId: user.UserId, FileName: path.Base(key), ContentType: contentType, MaxSize: info.Size, Size: info.Size, State: "ready", CompletedAt: &info.LastModified, Public: true, CreatedAt: time.Now()}
			var existing int64
			if err := tx.Model(&system.SysFile{}).Where("object_key = ?", next).Count(&existing).Error; err != nil {
				return err
			}
			if existing == 0 {
				if err := storage.RegisterExisting(tx, &metadata); err != nil {
					return err
				}
			}
			if err := tx.Where("object_key = ?", next).First(&metadata).Error; err != nil {
				return err
			}
			if err := tx.Unscoped().Model(&system.SysUser{}).Where("avatar = ?", user.Avatar).UpdateColumn("avatar", "/"+bucket+"/"+next).Error; err != nil {
				return err
			}
		}
		owners := make(map[string]bool, len(users))
		for _, user := range users {
			owners[user.TenantId+"/"+user.UserId] = true
		}
		// Rebuild ownership metadata for the legacy server-generated namespace.
		// Unknown objects stay private and require an explicit ownership decision.
		for object := range client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true}) {
			if object.Err != nil {
				return object.Err
			}
			parts := strings.Split(object.Key, "/")
			if len(parts) != 3 || !owners[parts[0]+"/"+parts[1]] {
				continue
			}
			info, err := client.StatObject(ctx, bucket, object.Key, minio.StatObjectOptions{})
			if err != nil {
				return err
			}
			if info.Size <= 0 {
				continue
			}
			contentType := info.ContentType
			if contentType == "" {
				contentType = strings.Split(mime.TypeByExtension(path.Ext(object.Key)), ";")[0]
			}
			if contentType == "" {
				contentType = "application/octet-stream"
			}
			metadata := system.SysFile{ObjectKey: object.Key, TenantId: parts[0], UserId: parts[1], FileName: path.Base(object.Key), ContentType: contentType, MaxSize: info.Size, Size: info.Size, State: "ready", CompletedAt: &info.LastModified, Public: false, CreatedAt: time.Now()}
			var existing int64
			if err := tx.Model(&system.SysFile{}).Where("object_key = ?", object.Key).Count(&existing).Error; err != nil {
				return err
			}
			if existing == 0 {
				if err := storage.RegisterExisting(tx, &metadata); err != nil {
					return err
				}
			}
			if err := tx.Where("object_key = ?", object.Key).First(&metadata).Error; err != nil {
				return err
			}
		}
		return tx.Exec("INSERT INTO sys_storage_migrations (name) VALUES (?)", "private_objects_v1").Error
	})
}
