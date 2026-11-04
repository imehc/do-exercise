package common

import (
	"errors"
	"fmt"
	"mime"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	request "github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/storage"
	"github.com/imehc/do-exercise/server/util"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type OssService struct{}

var allowedUploadExt = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true, ".ico": true, ".pdf": true, ".txt": true, ".csv": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true, ".zip": true}
var avatarExt = map[string]string{".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png", ".gif": "image/gif", ".webp": "image/webp"}

const signatureLifetime = storage.SignatureLifetime

func (s *OssService) GetPresignedUrl(db *gorm.DB, userID, tenantID string, req request.OssReq) (*request.OssRes, error) {
	if _, err := util.LoadSessionPrincipal(db, userID, tenantID); err != nil {
		return nil, err
	}
	ext := strings.ToLower(filepath.Ext(req.FileName))
	if !allowedUploadExt[ext] {
		return nil, errors.New("invalidFileType")
	}
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	contentType = strings.Split(contentType, ";")[0]
	public := req.Purpose == "avatar"
	prefix := "private"
	maxSize := int64(100 << 20)
	if public {
		var ok bool
		contentType, ok = avatarExt[ext]
		if !ok {
			return nil, errors.New("invalidFileType")
		}
		prefix = "public/avatars"
		maxSize = 10 << 20
	}
	if req.FileSize <= 0 || req.FileSize > maxSize {
		return nil, errors.New("invalidFileSize")
	}
	maxSize = req.FileSize
	id, err := util.Uuid()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s/%s/%s/%s%s", prefix, tenantID, userID, id, ext)
	policy := minio.NewPostPolicy()
	for _, err := range []error{policy.SetBucket(global.Config.Oss.BucketName), policy.SetKey(key), policy.SetExpires(time.Now().Add(signatureLifetime)), policy.SetContentType(contentType), policy.SetContentLengthRange(maxSize, maxSize)} {
		if err != nil {
			return nil, err
		}
	}
	target, fields, err := global.Oss.PresignedPostPolicy(db.Statement.Context, policy)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().UTC().Add(signatureLifetime)
	file := &system.SysFile{State: "pending", UploadExpiresAt: &expiresAt, ObjectKey: key, TenantId: tenantID, UserId: userID, FileName: filepath.Base(req.FileName), ContentType: contentType, MaxSize: maxSize, Public: public}
	if err := db.Transaction(func(tx *gorm.DB) error { return storage.Reserve(tx, file) }); err != nil {
		return nil, err
	}
	download := ""
	if public {
		download = "/" + global.Config.Oss.BucketName + "/" + key
	}
	return &request.OssRes{UploadUrl: externalObjectURL(target), FormFields: fields, ObjectKey: key, GetObjectUrl: download, Expires: int(signatureLifetime.Seconds())}, nil
}
func (s *OssService) GetDownloadUrl(db *gorm.DB, key string) (*request.OssDownloadRes, error) {
	var file system.SysFile
	if err := db.Where("object_key = ? AND state IN ?", key, []string{"ready", "legacy"}).First(&file).Error; err != nil {
		return nil, errors.New("fileNotFound")
	}
	object, err := global.Oss.StatObject(db.Statement.Context, global.Config.Oss.BucketName, key, minio.StatObjectOptions{})
	if err != nil || object.Size <= 0 || object.Size > file.MaxSize || object.ContentType != file.ContentType {
		return nil, errors.New("fileNotFound")
	}
	params := url.Values{"response-content-type": {file.ContentType}, "response-content-disposition": {"attachment"}}
	signed, err := global.Oss.PresignedGetObject(db.Statement.Context, global.Config.Oss.BucketName, key, signatureLifetime, params)
	if err != nil {
		return nil, err
	}
	return &request.OssDownloadRes{Url: externalObjectURL(signed), Expires: int(signatureLifetime.Seconds())}, nil
}
func externalObjectURL(original *url.URL) string {
	prefix := strings.TrimRight(global.Config.Oss.PresignedHost, "/")
	if gin.Mode() != gin.ReleaseMode || prefix == "" {
		return original.String()
	}
	if strings.HasPrefix(prefix, "/") {
		return prefix + original.RequestURI()
	}
	target, err := url.Parse(prefix)
	if err != nil || target.Host == "" {
		return original.String()
	}
	target.Path = strings.TrimRight(target.Path, "/") + original.Path
	target.RawQuery = original.RawQuery
	return target.String()
}
