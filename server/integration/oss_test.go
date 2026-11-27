//go:build integration

package integration

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/imehc/do-exercise/server/config"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/internal"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	service "github.com/imehc/do-exercise/server/service/common"
	"github.com/imehc/do-exercise/server/storage"
	"github.com/minio/minio-go/v7"
)

func TestObjectStorageTenantIsolation(t *testing.T) {
	endpoint := os.Getenv("OSS_TEST_ENDPOINT")
	if endpoint == "" {
		t.Skip("set OSS_TEST_ENDPOINT for real object-storage regression")
	}
	client, err := minio.New(endpoint, &minio.Options{Creds: objectStorageTestCredentials(), Region: "us-east-1", BucketLookup: minio.BucketLookupPath})
	must(t, err)
	bucket := "tenant-security-test"
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	must(t, err)
	if !exists {
		must(t, client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}))
	}
	must(t, client.SetBucketPolicy(ctx, bucket, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::tenant-security-test/public/avatars/*"]}]}`))
	oldOss, oldConfig := global.Oss, global.Config.Oss
	global.Oss = client
	global.Config.Oss = config.Oss{BucketName: bucket}
	t.Cleanup(func() { global.Oss = oldOss; global.Config.Oss = oldConfig })
	a, b := principal(t), principal(t)
	svc := service.OssService{}
	private, err := svc.GetPresignedUrl(tenantDB(a.TenantId), a.UserId, a.TenantId, common.OssReq{FileName: "sample.txt", FileSize: 19})
	must(t, err)
	avatar, err := svc.GetPresignedUrl(tenantDB(a.TenantId), a.UserId, a.TenantId, common.OssReq{FileName: "sample.png", FileSize: 11, Purpose: "avatar"})
	must(t, err)
	for _, item := range []*common.OssRes{private, avatar} {
		t.Cleanup(func() { _ = client.RemoveObject(ctx, bucket, item.ObjectKey, minio.RemoveObjectOptions{}) })
	}
	// upload 返回对象存储的响应状态码。err 非 nil 表示连接在拿到响应之前就中断了：
	// 对“越界/类型不符”这类会被服务端提前拒绝的请求，连接中断本身就等于被拒绝，
	// 由调用方按语义判断，不能一律当成测试失败。
	upload := func(item *common.OssRes, content []byte, tamper string) (int, error) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		for key, value := range item.FormFields {
			if tamper != "" && key == "Content-Type" {
				value = tamper
			}
			must(t, writer.WriteField(key, value))
		}
		part, err := writer.CreateFormFile("file", "sample")
		must(t, err)
		_, err = part.Write(content)
		must(t, err)
		must(t, writer.Close())
		request, err := http.NewRequest(http.MethodPost, item.UploadUrl, &body)
		must(t, err)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return 0, err
		}
		defer response.Body.Close()
		_, _ = io.Copy(io.Discard, response.Body)
		return response.StatusCode, nil
	}
	if status, err := upload(private, []byte("private tenant data"), ""); err != nil || status < 200 || status >= 300 {
		t.Fatalf("signed private upload: status=%d err=%v", status, err)
	}
	if status, err := upload(avatar, []byte("image bytes"), ""); err != nil || status < 200 || status >= 300 {
		t.Fatalf("signed avatar upload: status=%d err=%v", status, err)
	}
	if status, err := upload(avatar, []byte("bad mime"), "text/html"); err != nil || status < 400 {
		t.Fatalf("tampered MIME accepted: status=%d err=%v", status, err)
	}
	// 预签名策略的 content-length-range 是 [FileSize, FileSize]（avatar 为 11 字节），
	// 只要超过 11 字节就会被同一条策略拒绝，因此探针用 12 字节即可，不必灌 10 MiB。
	// 服务端会在读完 body 之前关闭连接，客户端侧表现为 broken pipe / connection reset；
	// 这种中断本身就是“请求被拒绝”，只有连接正常结束且返回 2xx 才是越界真的被接受。
	if status, err := upload(avatar, make([]byte, 12), ""); err == nil && status < 400 {
		t.Fatalf("oversize accepted: %d", status)
	}
	get := func(target string) int {
		response, err := http.Get(target)
		must(t, err)
		defer response.Body.Close()
		_, _ = io.Copy(io.Discard, response.Body)
		return response.StatusCode
	}
	origin := "http://" + endpoint + "/" + bucket + "/"
	if status := get(origin + private.ObjectKey); status != 403 {
		t.Fatalf("anonymous private GET: %d", status)
	}
	if status := get(origin + "?list-type=2"); status != 403 {
		t.Fatalf("anonymous list: %d", status)
	}
	if status := get(origin + avatar.ObjectKey); status != 200 {
		t.Fatalf("public avatar GET: %d", status)
	}
	if _, err := svc.GetDownloadUrl(tenantDB(b.TenantId), private.ObjectKey); err == nil {
		t.Fatal("foreign download authorized")
	}
	_, err = storage.Complete(tenantDB(a.TenantId), a.UserId, private.ObjectKey)
	must(t, err)
	signed, err := svc.GetDownloadUrl(tenantDB(a.TenantId), private.ObjectKey)
	must(t, err)
	if status := get(signed.Url); status != 200 {
		t.Fatalf("authorized download: %d", status)
	}
	// Existing avatar references are copied exactly once; unrelated files remain private.
	legacy := a.TenantId + "/" + a.UserId + "/old.png"
	_, err = client.PutObject(ctx, bucket, legacy, bytes.NewReader([]byte("old avatar")), 10, minio.PutObjectOptions{ContentType: "image/png"})
	must(t, err)
	must(t, bypass().Model(&system.SysUser{}).Where("id = ?", a.UserId).UpdateColumn("avatar", "/"+bucket+"/"+legacy).Error)
	must(t, bypass().Exec("DROP TABLE IF EXISTS sys_storage_migrations").Error)
	must(t, internal.MigrateLegacyAvatars(ctx, client))
	var migrated system.SysUser
	must(t, bypass().First(&migrated, "id = ?", a.UserId).Error)
	if migrated.Avatar == "/"+bucket+"/"+legacy {
		t.Fatal("legacy avatar not migrated")
	}
	if status := get("http://" + endpoint + migrated.Avatar); status != 200 {
		t.Fatalf("migrated avatar GET: %d", status)
	}
	must(t, internal.MigrateLegacyAvatars(ctx, client))
	legacyDownload, err := svc.GetDownloadUrl(tenantDB(a.TenantId), legacy)
	must(t, err)
	if status := get(legacyDownload.Url); status != 200 {
		t.Fatalf("legacy private file inaccessible: %d", status)
	}
	if status := get(origin + legacy); status != 403 {
		t.Fatalf("legacy namespace still public: %d", status)
	}
}
