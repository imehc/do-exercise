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
	upload := func(item *common.OssRes, content []byte, tamper string) int {
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
		must(t, err)
		defer response.Body.Close()
		_, _ = io.Copy(io.Discard, response.Body)
		return response.StatusCode
	}
	if status := upload(private, []byte("private tenant data"), ""); status < 200 || status >= 300 {
		t.Fatalf("signed private upload: %d", status)
	}
	if status := upload(avatar, []byte("image bytes"), ""); status < 200 || status >= 300 {
		t.Fatalf("signed avatar upload: %d", status)
	}
	if status := upload(avatar, []byte("bad mime"), "text/html"); status < 400 {
		t.Fatalf("tampered MIME accepted: %d", status)
	}
	if status := upload(avatar, make([]byte, (10<<20)+1), ""); status < 400 {
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
