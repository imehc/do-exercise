package internal

import (
	"context"
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/util"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// InitOss 初始化RustFS
func InitOss() {
	cfg := global.Config.Oss
	client, err := minio.New(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.Secure, // 是否使用https进行通信，由配置驱动
			// 单区域自建部署固定 us-east-1：显式指定后 minio-go 不再发起
			// GetBucketLocation 探测（RustFS 对 ?location 的实现差异不再影响签名作用域）
			Region: "us-east-1",
			// 路径式寻址：虚拟主机式(bucket.host)在自建 + nginx 反代下不可用
			BucketLookup: minio.BucketLookupPath,
		},
	)
	if err != nil {
		util.Exit("oss连接失败", err)
	}

	// 检查bucket是否存在
	bucketName := cfg.BucketName
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		util.Exit("oss检查bucket失败", err)
	}
	if !exists {
		// 创建bucket
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			util.Exit("oss创建bucket失败", err)
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
	fmt.Println("oss初始化成功")
}
