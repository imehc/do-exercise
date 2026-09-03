package global

import (
	"context"
	"net/url"
	"time"
)

// OssClient 对象存储客户端的最小抽象。
//
// 目前运行期只用预签名上传（Put），因此接口只暴露 PresignedPutObject；
// 需要签名下载/对象清理等方法时（对象存储访问控制改造）再往接口补，
// 不提前为未消费的能力定义签名。
type OssClient interface {
	PresignedPutObject(ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error)
}
