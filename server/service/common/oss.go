package common

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	model "github.com/imehc/do-exercise/server/model/common"
)

type OssService struct{}

// GetPresignedUrl 获取预签名
func (s *OssService) GetPresignedUrl(req model.OssReq) (*model.OssRes, error) {
	client := global.Oss
	bucketName := global.Config.Minio.BucketName
	expires := time.Duration(global.Config.Minio.Expires) * time.Second
	ctx := context.Background()
	putUrl, err := client.PresignedPutObject(
		ctx,
		bucketName,
		req.FileName,
		expires,
	)
	if err != nil {
		return nil, errors.New("getFailed")
	}

	// getUrl, err := client.PresignedGetObject(
	// 	ctx,
	// 	bucketName,
	// 	req.FileName,
	// 	expires,
	// 	nil,
	// )
	// if err != nil {
	// 	return nil, errors.New("getFailed")
	// }

	url := fmt.Sprintf("%s://%s%s", putUrl.Scheme, putUrl.Host, putUrl.Path)
	presignedHost := strings.TrimRight(global.Config.Minio.PresignedHost, "/")
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
		Expires:      global.Config.Minio.Expires,
	}, nil
}
