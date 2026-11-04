//go:build integration

package integration

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/imehc/do-exercise/server/config"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	service "github.com/imehc/do-exercise/server/service/common"
	"github.com/imehc/do-exercise/server/storage"
	"github.com/minio/minio-go/v7"
)

func storageClient(t *testing.T) *minio.Client {
	t.Helper()
	endpoint := os.Getenv("OSS_TEST_ENDPOINT")
	if endpoint == "" {
		t.Skip("OSS_TEST_ENDPOINT is required")
	}
	client, err := minio.New(endpoint, &minio.Options{Creds: objectStorageTestCredentials(), Region: "us-east-1", BucketLookup: minio.BucketLookupPath})
	must(t, err)
	bucket := "tenant-lifecycle-test"
	exists, err := client.BucketExists(context.Background(), bucket)
	must(t, err)
	if !exists {
		must(t, client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{}))
	}
	old, oldConfig := global.Oss, global.Config.Oss
	global.Oss = client
	global.Config.Oss = config.Oss{BucketName: bucket}
	t.Cleanup(func() { global.Oss = old; global.Config.Oss = oldConfig })
	return client
}

type failedDeleteClient struct{ global.OssClient }

func (failedDeleteClient) RemoveObject(context.Context, string, string, minio.RemoveObjectOptions) error {
	return errors.New("storage unavailable")
}

func TestStorageReservationsAndCollection(t *testing.T) {
	truncateScopedTables(t)
	client := storageClient(t)
	u, other := principal(t), principal(t)
	must(t, bypass().Model(&system.SysTenant{}).Where("tenant_id = ?", u.TenantId).UpdateColumn("max_storage_bytes", 100).Error)
	svc := service.OssService{}
	var success atomic.Int32
	var mu sync.Mutex
	var keys []string
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			signed, err := svc.GetPresignedUrl(tenantDB(u.TenantId), u.UserId, u.TenantId, common.OssReq{FileName: "quota.txt", FileSize: 10})
			if err == nil {
				success.Add(1)
				mu.Lock()
				keys = append(keys, signed.ObjectKey)
				mu.Unlock()
			} else if err.Error() != "tenantStorageQuotaReached" {
				t.Errorf("reservation: %v", err)
			}
		}()
	}
	wg.Wait()
	if success.Load() != 10 {
		t.Fatalf("quota admitted %d uploads", success.Load())
	}
	charged, limit, err := storage.Usage(tenantDB(u.TenantId), u.TenantId)
	must(t, err)
	if charged != 100 || limit != 100 {
		t.Fatalf("usage %d/%d", charged, limit)
	}
	key := keys[0]
	if _, err := storage.Complete(tenantDB(u.TenantId), u.UserId, key); err == nil {
		t.Fatal("unuploaded file completed")
	}
	_, err = client.PutObject(context.Background(), global.Config.Oss.BucketName, key, bytes.NewReader([]byte("0123456789")), 10, minio.PutObjectOptions{ContentType: "text/plain"})
	must(t, err)
	if _, err := storage.Complete(tenantDB(other.TenantId), other.UserId, key); err == nil {
		t.Fatal("foreign completion accepted")
	}
	_, err = storage.Complete(tenantDB(u.TenantId), u.UserId, key)
	must(t, err)
	_, err = storage.Complete(tenantDB(u.TenantId), u.UserId, key)
	must(t, err)
	if err := storage.Delete(tenantDB(other.TenantId), other.UserId, key); err == nil {
		t.Fatal("foreign deletion accepted")
	}
	must(t, storage.Delete(tenantDB(u.TenantId), u.UserId, key))
	now := time.Now().UTC()
	n, err := storage.Collect(context.Background(), db, now, 1000)
	must(t, err)
	if n != 0 {
		t.Fatalf("collected %d files before signature expiry", n)
	}
	global.Oss = failedDeleteClient{client}
	if _, err := storage.Collect(context.Background(), db, now.Add(11*time.Minute), 1000); err == nil {
		t.Fatal("delete error lost")
	}
	charged, _, err = storage.Usage(tenantDB(u.TenantId), u.TenantId)
	must(t, err)
	if charged != 100 {
		t.Fatal("failed deletion released quota")
	}
	global.Oss = client
	_, err = storage.Collect(context.Background(), db, now.Add(11*time.Minute), 1000)
	must(t, err)
	charged, _, err = storage.Usage(tenantDB(u.TenantId), u.TenantId)
	must(t, err)
	if charged != 0 {
		t.Fatalf("quota leaked: %d", charged)
	}
	if _, err := client.StatObject(context.Background(), global.Config.Oss.BucketName, key, minio.StatObjectOptions{}); err == nil {
		t.Fatal("deleted object remains")
	}
	// An active private file has no age-based expiry; avatars are reference-aware.
	for _, purpose := range []string{"private", "avatar"} {
		ext := ".txt"
		mime := "text/plain"
		if purpose == "avatar" {
			ext = ".png"
			mime = "image/png"
		}
		signed, err := svc.GetPresignedUrl(tenantDB(u.TenantId), u.UserId, u.TenantId, common.OssReq{FileName: "live" + ext, FileSize: 10, Purpose: purpose})
		must(t, err)
		_, err = client.PutObject(context.Background(), global.Config.Oss.BucketName, signed.ObjectKey, bytes.NewReader([]byte("0123456789")), 10, minio.PutObjectOptions{ContentType: mime})
		must(t, err)
		_, err = storage.Complete(tenantDB(u.TenantId), u.UserId, signed.ObjectKey)
		must(t, err)
		if purpose == "avatar" {
			must(t, bypass().Model(&system.SysUser{}).Where("id = ?", u.UserId).UpdateColumn("avatar", "/"+global.Config.Oss.BucketName+"/"+signed.ObjectKey).Error)
			if err := storage.Delete(tenantDB(u.TenantId), u.UserId, signed.ObjectKey); err == nil {
				t.Fatal("referenced avatar deleted")
			}
		}
		_, err = storage.Collect(context.Background(), db, now.Add(48*time.Hour), 1000)
		must(t, err)
		_, err = client.StatObject(context.Background(), global.Config.Oss.BucketName, signed.ObjectKey, minio.StatObjectOptions{})
		must(t, err)
		if purpose == "avatar" {
			must(t, bypass().Model(&system.SysUser{}).Where("id = ?", u.UserId).UpdateColumn("avatar", "").Error)
			_, err = storage.Collect(context.Background(), db, now.Add(48*time.Hour), 1000)
			must(t, err)
			if _, err := client.StatObject(context.Background(), global.Config.Oss.BucketName, signed.ObjectKey, minio.StatObjectOptions{}); err == nil {
				t.Fatal("unused avatar not collected")
			}
		}
	}
}

func jobFixture(t *testing.T, concurrent bool, retries, timeout uint) *system.SysJob {
	t.Helper()
	u := principal(t)
	job := &system.SysJob{Name: "execution", JobGroup: "g", TenantId: u.TenantId, CronExpression: "0 0 0 * * *", Command: "clean_empty_username_operation_logs", Status: 1, Concurrent: concurrent, RetryTimes: retries, RetryInterval: 1, Timeout: timeout}
	must(t, bypass().Create(job).Error)
	return job
}
func installJobTrigger(t *testing.T, body string) {
	t.Helper()
	must(t, bypass().Exec("CREATE OR REPLACE FUNCTION preprod_job_trigger() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN "+body+" RETURN NULL; END $$").Error)
	must(t, bypass().Exec("CREATE TRIGGER preprod_job_trigger BEFORE DELETE ON sys_operation_log FOR EACH STATEMENT EXECUTE FUNCTION preprod_job_trigger()").Error)
	t.Cleanup(func() {
		_ = bypass().Exec("DROP TRIGGER IF EXISTS preprod_job_trigger ON sys_operation_log").Error
		_ = bypass().Exec("DROP FUNCTION IF EXISTS preprod_job_trigger()").Error
		_ = bypass().Exec("DROP SEQUENCE IF EXISTS preprod_job_attempts").Error
	})
}
func waitRun(t *testing.T, jobID uint, status string) system.SysJobRun {
	t.Helper()
	deadline := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		var run system.SysJobRun
		must(t, bypass().Where("job_id = ?", jobID).Order("created_at DESC").First(&run).Error)
		if run.Status == status {
			return run
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("job %d did not reach %s", jobID, status)
	return system.SysJobRun{}
}
func TestJobRetryRollbackAndRecovery(t *testing.T) {
	must(t, bypass().Exec("TRUNCATE sys_job_run").Error)
	js := shared.GetJobScheduler()
	js.Stop()
	t.Cleanup(js.Stop)
	job := jobFixture(t, false, 2, 5)
	must(t, bypass().Exec("CREATE SEQUENCE preprod_job_attempts").Error)
	installJobTrigger(t, "IF nextval('preprod_job_attempts') <= 2 THEN RAISE EXCEPTION 'transient job failure'; END IF;")
	log := &system.SysOperationLog{TenantId: job.TenantId}
	must(t, bypass().Create(log).Error)
	if err := js.ExecuteJob(job); err == nil {
		t.Fatal("first attempt unexpectedly succeeded")
	}
	retry := waitRun(t, job.Id, "retrying")
	if retry.Attempts != 1 {
		t.Fatalf("attempts %d", retry.Attempts)
	}
	var count int64
	must(t, tenantDB(job.TenantId).Model(&system.SysOperationLog{}).Count(&count).Error)
	if count != 1 {
		t.Fatal("failed handler effects committed")
	}
	// Another scheduler invocation can consume the persisted retry after restart.
	js.Start()
	run := waitRun(t, job.Id, "succeeded")
	js.Stop()
	if run.Attempts != 3 {
		t.Fatalf("retries %d", run.Attempts)
	}
	var got system.SysJob
	must(t, bypass().First(&got, job.Id).Error)
	if got.Times != 1 {
		t.Fatalf("success count %d", got.Times)
	}
	must(t, tenantDB(job.TenantId).Model(&system.SysOperationLog{}).Count(&count).Error)
	if count != 0 {
		t.Fatal("successful effects missing")
	}
}
func TestJobTimeoutAndRetryExhaustion(t *testing.T) {
	js := shared.GetJobScheduler()
	js.Stop()
	t.Cleanup(js.Stop)
	job := jobFixture(t, false, 1, 1)
	installJobTrigger(t, "PERFORM pg_sleep(2);")
	if err := js.ExecuteJob(job); err == nil {
		t.Fatal("slow handler accepted")
	}
	first := waitRun(t, job.Id, "retrying")
	if first.Attempts != 1 {
		t.Fatal("timeout was not persisted")
	}
	js.Start()
	run := waitRun(t, job.Id, "failed")
	js.Stop()
	if run.Attempts != 2 {
		t.Fatalf("timeout attempts %d", run.Attempts)
	}
	var got system.SysJob
	must(t, bypass().First(&got, job.Id).Error)
	if got.Times != 0 {
		t.Fatal("timeout counted as success")
	}
}
func TestJobConcurrentAdmissionAndReplicaDedup(t *testing.T) {
	js := &shared.JobScheduler{}
	serial := jobFixture(t, false, 0, 5)
	run, err := js.SubmitJob(serial)
	must(t, err)
	if run == nil {
		t.Fatal("not queued")
	}
	if _, err := js.SubmitJob(serial); err == nil {
		t.Fatal("serial job admitted overlap")
	}
	parallel := jobFixture(t, true, 0, 5)
	for i := 0; i < 4; i++ {
		_, err := js.SubmitJob(parallel)
		must(t, err)
	}
	if _, err := js.SubmitJob(parallel); err == nil {
		t.Fatal("unbounded concurrency")
	}
	// Multiple replicas racing the same scheduled occurrence admit only one run.
	dedup := jobFixture(t, true, 0, 5)
	slot := time.Now().UTC().Truncate(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := js.ExecuteScheduledJob(dedup, slot); err != nil {
				t.Errorf("occurrence: %v", err)
			}
		}()
	}
	wg.Wait()
	var count int64
	must(t, bypass().Model(&system.SysJobRun{}).Where("job_id = ?", dedup.Id).Count(&count).Error)
	if count != 1 {
		t.Fatalf("duplicate occurrences %d", count)
	}
	// SQL statement work can overlap; there is no job lock around the handler.
	installJobTrigger(t, "PERFORM pg_sleep(0.4);")
	parallel = jobFixture(t, true, 0, 5)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := js.ExecuteJob(parallel); err != nil {
				t.Errorf("parallel execution: %v", err)
			}
		}(i)
	}
	// Observe overlapping SQL statements instead of relying on a narrow latency threshold.
	observed := int64(0)
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		var active int64
		must(t, bypass().Raw(`SELECT count(*) FROM pg_stat_activity WHERE pid <> pg_backend_pid() AND state = 'active' AND query LIKE 'DELETE FROM "sys_operation_log"%'`).Scan(&active).Error)
		if active > observed {
			observed = active
		}
		if observed >= 2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
	if observed < 2 {
		t.Fatal("database handlers did not overlap")
	}

}

func TestUploadPolicyRejectsSlowExpiredPost(t *testing.T) {
	proxy := os.Getenv("OSS_TEST_UPLOAD_PROXY")
	if proxy == "" {
		t.Skip("OSS_TEST_UPLOAD_PROXY must point to the production upload ingress")
	}
	client := storageClient(t)
	key := "private/slow-policy/expired.txt"
	policy := minio.NewPostPolicy()
	must(t, policy.SetBucket(global.Config.Oss.BucketName))
	must(t, policy.SetKey(key))
	must(t, policy.SetExpires(time.Now().Add(time.Second)))
	must(t, policy.SetContentType("text/plain"))
	must(t, policy.SetContentLengthRange(10, 10))
	target, fields, err := client.PresignedPostPolicy(context.Background(), policy)
	must(t, err)
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	for k, v := range fields {
		must(t, form.WriteField(k, v))
	}
	part, err := form.CreateFormFile("file", "slow.txt")
	must(t, err)
	_, err = part.Write([]byte("0123456789"))
	must(t, err)
	must(t, form.Close())
	reader, writer := io.Pipe()
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer writer.Close()
		data := body.Bytes()
		_, _ = writer.Write(data[:bytes.Index(data, []byte("0123456789"))+5])
		time.Sleep(3 * time.Second)
		_, _ = writer.Write(data[bytes.Index(data, []byte("0123456789"))+5:])
	}()
	request, err := http.NewRequest(http.MethodPost, strings.TrimRight(proxy, "/")+target.Path, reader)
	must(t, err)
	request.ContentLength = int64(body.Len())
	request.Header.Set("Content-Type", form.FormDataContentType())
	response, err := (&http.Client{Timeout: 10 * time.Second}).Do(request)
	must(t, err)
	_, _ = io.Copy(io.Discard, response.Body)
	response.Body.Close()
	reader.Close()
	<-done
	t.Cleanup(func() {
		_ = client.RemoveObject(context.Background(), global.Config.Oss.BucketName, key, minio.RemoveObjectOptions{})
	})
	if response.StatusCode < 400 {
		t.Fatalf("upload completed after policy expiry: %d", response.StatusCode)
	}
	if _, err := client.StatObject(context.Background(), global.Config.Oss.BucketName, key, minio.StatObjectOptions{}); err == nil {
		t.Fatal("expired upload left an object")
	}
}

func TestMaximumUploadThroughIngress(t *testing.T) {
	proxy := os.Getenv("OSS_TEST_UPLOAD_PROXY")
	if proxy == "" {
		t.Skip("OSS_TEST_UPLOAD_PROXY required")
	}
	client := storageClient(t)
	user := principal(t)
	const size = 100 << 20
	signed, err := (&service.OssService{}).GetPresignedUrl(tenantDB(user.TenantId), user.UserId, user.TenantId, common.OssReq{FileName: "boundary.zip", FileSize: size})
	must(t, err)
	target, err := url.Parse(signed.UploadUrl)
	must(t, err)
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	for k, v := range signed.FormFields {
		must(t, form.WriteField(k, v))
	}
	part, err := form.CreateFormFile("file", "boundary.zip")
	must(t, err)
	_, err = io.CopyN(part, repeatZeroReader{}, size)
	must(t, err)
	must(t, form.Close())
	request, err := http.NewRequest(http.MethodPost, strings.TrimRight(proxy, "/")+target.Path, &body)
	must(t, err)
	request.Header.Set("Content-Type", form.FormDataContentType())
	started := time.Now()
	response, err := (&http.Client{Timeout: 2 * time.Minute}).Do(request)
	must(t, err)
	_, _ = io.Copy(io.Discard, response.Body)
	response.Body.Close()
	t.Cleanup(func() {
		_ = client.RemoveObject(context.Background(), global.Config.Oss.BucketName, signed.ObjectKey, minio.RemoveObjectOptions{})
	})
	if response.StatusCode >= 300 {
		t.Fatalf("maximum upload rejected: %d", response.StatusCode)
	}
	file, err := storage.Complete(tenantDB(user.TenantId), user.UserId, signed.ObjectKey)
	must(t, err)
	if file.Size != size {
		t.Fatalf("actual size %d", file.Size)
	}
	charged, _, err := storage.Usage(tenantDB(user.TenantId), user.TenantId)
	must(t, err)
	if charged != size {
		t.Fatalf("quota size %d", charged)
	}
	t.Logf("Ingress upload size=100MiB duration=%s quota=100MiB", time.Since(started))
}

type repeatZeroReader struct{}

func (repeatZeroReader) Read(p []byte) (int, error) { clear(p); return len(p), nil }

func TestTenantStorageRetention(t *testing.T) {
	client := storageClient(t)
	u := principal(t)
	signed, err := (&service.OssService{}).GetPresignedUrl(tenantDB(u.TenantId), u.UserId, u.TenantId, common.OssReq{FileName: "retained.txt", FileSize: 10})
	must(t, err)
	_, err = client.PutObject(context.Background(), global.Config.Oss.BucketName, signed.ObjectKey, bytes.NewReader([]byte("0123456789")), 10, minio.PutObjectOptions{ContentType: "text/plain"})
	must(t, err)
	_, err = storage.Complete(tenantDB(u.TenantId), u.UserId, signed.ObjectKey)
	must(t, err)
	now := time.Now().UTC()
	expired := now.Add(-24 * time.Hour)
	must(t, bypass().Model(&system.SysTenant{}).Where("tenant_id = ?", u.TenantId).UpdateColumns(map[string]any{"status": false, "expire_time": expired}).Error)
	_, err = storage.Collect(context.Background(), db, now.Add(11*time.Minute), 1000)
	must(t, err)
	_, err = client.StatObject(context.Background(), global.Config.Oss.BucketName, signed.ObjectKey, minio.StatObjectOptions{})
	must(t, err)
	must(t, bypass().Model(&system.SysTenant{}).Where("tenant_id = ?", u.TenantId).UpdateColumn("deleted_at", now.Add(-29*24*time.Hour)).Error)
	_, err = storage.Collect(context.Background(), db, now.Add(11*time.Minute), 1000)
	must(t, err)
	_, err = client.StatObject(context.Background(), global.Config.Oss.BucketName, signed.ObjectKey, minio.StatObjectOptions{})
	must(t, err)
	must(t, bypass().Unscoped().Model(&system.SysTenant{}).Where("tenant_id = ?", u.TenantId).UpdateColumn("deleted_at", now.Add(-31*24*time.Hour)).Error)
	_, err = storage.Collect(context.Background(), db, now.Add(11*time.Minute), 1000)
	must(t, err)
	if _, err := client.StatObject(context.Background(), global.Config.Oss.BucketName, signed.ObjectKey, minio.StatObjectOptions{}); err == nil {
		t.Fatal("deleted tenant object not collected")
	}
	var usage system.SysStorageUsage
	must(t, tenantDB(u.TenantId).First(&usage, "tenant_id = ?", u.TenantId).Error)
	if usage.ChargedBytes != 0 {
		t.Fatal("deleted tenant quota not released")
	}
}
