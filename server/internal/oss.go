package internal

import (
	"context"
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/util"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// InitMinio 初始化Minio
func InitMinio() {
	cfg := global.Config.Minio
	client, err := minio.New(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: false, // 是否使用https进行通信
		},
	)
	if err != nil {
		util.Exit("minio连接失败", err)
	}

	// 检查bucket是否存在
	bucketName := cfg.BucketName
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		util.Exit("minio检查bucket失败", err)
	}
	if !exists {
		// 创建bucket
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			util.Exit("minio创建bucket失败", err)
		}
	}

	// 设置匿名访问规则
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"AWS": ["*"]
				},
				"Action": [
					"s3:GetBucketLocation",
					"s3:ListBucket"
				],
				"Resource": ["arn:aws:s3:::%s"]
			},
			{
				"Effect": "Allow",
				"Principal": {
					"AWS": ["*"]
				},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName, bucketName)
	err = client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		util.Exit("设置匿名访问规则失败", err)
	}

	global.Oss = client
	fmt.Println("minio初始化成功")
}
