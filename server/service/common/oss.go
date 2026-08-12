package common

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	model "github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/util"
)

type OssService struct{}

// allowedUploadExt 允许上传的文件扩展名。
//
// 刻意排除 .svg / .html / .htm / .xhtml：bucket 策略允许匿名 GetObject
// （见 internal/oss.go），这类文件一旦上传就会成为你域名下的永久公开页面，
// 可直接用于钓鱼或存储型 XSS。
var allowedUploadExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
	".ico":  true,
	".pdf":  true,
	".txt":  true,
	".csv":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
	".zip":  true,
}

// GetPresignedUrl 获取预签名上传地址。
//
// 对象 key 一律由服务端生成：预签名 PUT 在设计上绕过服务端校验，
// 客户端若能决定 key，就能覆盖任意已有对象（例如他人头像），
// 因此 file_name 只用来取扩展名，不参与 key 的构造。
func (s *OssService) GetPresignedUrl(userId string, req model.OssReq) (*model.OssRes, error) {
	if userId == "" {
		return nil, errors.New("getFailed")
	}

	ext := strings.ToLower(filepath.Ext(req.FileName))
	if !allowedUploadExt[ext] {
		return nil, errors.New("invalidFileType")
	}

	objectName, err := util.Uuid()
	if err != nil {
		return nil, errors.New("getFailed")
	}
	// 以 userId 作前缀，便于配额统计与按用户清理
	objectKey := fmt.Sprintf("%s/%s%s", userId, objectName, ext)

	client := global.Oss
	bucketName := global.Config.Oss.BucketName
	expires := time.Duration(global.Config.Oss.Expires) * time.Second
	ctx := context.Background()
	putUrl, err := client.PresignedPutObject(
		ctx,
		bucketName,
		objectKey,
		expires,
	)
	if err != nil {
		return nil, errors.New("getFailed")
	}
	url := fmt.Sprintf("%s://%s%s", putUrl.Scheme, putUrl.Host, putUrl.Path)
	presignedHost := strings.TrimRight(global.Config.Oss.PresignedHost, "/")
	if gin.Mode() == gin.ReleaseMode && presignedHost != "" {
		if strings.HasPrefix(presignedHost, "/") {
			url = fmt.Sprintf("%s%s", presignedHost, putUrl.Path)
		} else {
			url = fmt.Sprintf("%s://%s%s", putUrl.Scheme, presignedHost, putUrl.Path)
		}
	}

	return &model.OssRes{
		PutObjectUrl: fmt.Sprintf("%s?%s", url, putUrl.RawQuery),
		GetObjectUrl: putUrl.Path,
		Expires:      global.Config.Oss.Expires,
	}, nil
}
