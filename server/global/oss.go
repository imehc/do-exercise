package global

import (
	"context"
	"github.com/minio/minio-go/v7"
	"net/url"
	"time"
)

type OssClient interface {
	PresignedPostPolicy(context.Context, *minio.PostPolicy) (*url.URL, map[string]string, error)
	PresignedGetObject(context.Context, string, string, time.Duration, url.Values) (*url.URL, error)
	RemoveObject(context.Context, string, string, minio.RemoveObjectOptions) error
	StatObject(context.Context, string, string, minio.StatObjectOptions) (minio.ObjectInfo, error)
}
