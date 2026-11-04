//go:build integration

package integration

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	svc "github.com/imehc/do-exercise/server/service/common"
	"github.com/imehc/do-exercise/server/storage"
	"github.com/imehc/do-exercise/server/util"
	"github.com/minio/minio-go/v7"
)

// Opt-in workloads are excluded from normal CI. Targets must point to isolated
// preproduction processes sharing this integration database and Redis DB 15.
func TestPreproductionLoad(t *testing.T) {
	if os.Getenv("PREPROD_RUN") != "1" {
		t.Skip("set PREPROD_RUN=1 and PREPROD_URLS")
	}
	targets := strings.Split(os.Getenv("PREPROD_URLS"), ",")
	if len(targets) < 2 {
		t.Fatal("at least two backend URLs required")
	}
	client := storageClient(t)
	users := make([]*system.SysUser, 10)
	tokens := make([]*common.Token, 10)
	for i := range users {
		u := principal(t)
		users[i] = u
		role := &system.SysRole{Name: "load", Code: "load", TenantId: u.TenantId}
		must(t, bypass().Create(role).Error)
		var menus []system.SysMenu
		must(t, bypass().Where("scope <> ?", global.MenuScopePlatform).Find(&menus).Error)
		must(t, bypass().Model(role).Association("Menus").Append(menus))
		must(t, bypass().Model(u).Association("Roles").Append(role))
		token, err := (&util.Token{UserId: u.UserId, TenantId: u.TenantId, ExpireTime: 10 * time.Minute, RefreshExpireTime: time.Hour}).GenerateToken()
		must(t, err)
		tokens[i] = token
	}
	httpClient := &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{MaxIdleConns: 128, MaxIdleConnsPerHost: 64}}
	paths := []string{"/user/profile", "/system/users?page=1&page_size=10", "/system/jobs?page=1&page_size=10", "/files/usage"}
	for _, workers := range []int{16, 64} {
		var wg sync.WaitGroup
		var completed, limited, failures atomic.Int64
		var latencyMu sync.Mutex
		var latencies []time.Duration
		started := time.Now()
		deadline := started.Add(20 * time.Second)
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func(w int) {
				defer wg.Done()
				counter := 0
				for time.Now().Before(deadline) {
					tenant := (w + counter) % len(users)
					path := paths[counter%len(paths)]
					target := targets[(w+counter)%len(targets)]
					request, err := http.NewRequest(http.MethodGet, target+path, nil)
					if err != nil {
						failures.Add(1)
						continue
					}
					request.Header.Set("Authorization", tokens[tenant].AccessToken)
					// Test ingress trusts loopback; use a population of client addresses.
					request.Header.Set("X-Forwarded-For", fmt.Sprintf("10.%d.%d.%d", w+1, (counter/200)%200+1, counter%200+1))
					begin := time.Now()
					response, err := httpClient.Do(request)
					elapsed := time.Since(begin)
					if err != nil {
						failures.Add(1)
						counter++
						continue
					}
					body, readErr := io.ReadAll(response.Body)
					response.Body.Close()
					switch response.StatusCode {
					case 200:
						completed.Add(1)
						latencyMu.Lock()
						latencies = append(latencies, elapsed)
						latencyMu.Unlock()
						if readErr != nil {
							failures.Add(1)
						}
						if path == "/user/profile" {
							var profile struct {
								ID string `json:"id"`
							}
							if json.Unmarshal(body, &profile) != nil || profile.ID != users[tenant].UserId {
								failures.Add(1)
							}
						}
					case 429:
						limited.Add(1)
					default:
						failures.Add(1)
					}
					counter++
				}
			}(w)
		}
		wg.Wait()
		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
		percentile := func(p float64) time.Duration {
			if len(latencies) == 0 {
				return 0
			}
			return latencies[int(float64(len(latencies)-1)*p)]
		}
		t.Logf("HTTP workers=%d tenants=10 replicas=2 duration=%s success=%d throughput=%.1f/s p50=%s p95=%s p99=%s rate_limited=%d errors=%d", workers, time.Since(started), completed.Load(), float64(completed.Load())/time.Since(started).Seconds(), percentile(.5), percentile(.95), percentile(.99), limited.Load(), failures.Load())
		if failures.Load() != 0 || completed.Load() == 0 {
			t.Fatal("HTTP workload failed")
		}
	}
	// Real signed uploads through both backend processes, including completion and quota.
	time.Sleep(1100 * time.Millisecond)
	var wg sync.WaitGroup
	var failures atomic.Int64
	started := time.Now()
	for w := 0; w < 16; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 8; i++ {
				tenant := (w + i) % 10
				target := targets[(w+i)%2]
				request, _ := http.NewRequest(http.MethodGet, target+"/get_presigned_url?file_name=load.txt&file_size=4096", nil)
				request.Header.Set("Authorization", tokens[tenant].AccessToken)
				request.Header.Set("X-Forwarded-For", fmt.Sprintf("11.%d.%d.1", w+1, i+1))
				response, err := httpClient.Do(request)
				if err != nil {
					failures.Add(1)
					continue
				}
				var signed common.OssRes
				err = json.NewDecoder(response.Body).Decode(&signed)
				response.Body.Close()
				if response.StatusCode != 200 || err != nil {
					failures.Add(1)
					continue
				}
				var body bytes.Buffer
				form := multipart.NewWriter(&body)
				for k, v := range signed.FormFields {
					_ = form.WriteField(k, v)
				}
				part, _ := form.CreateFormFile("file", "load.txt")
				_, _ = part.Write(bytes.Repeat([]byte("x"), 4096))
				_ = form.Close()
				post, _ := http.NewRequest(http.MethodPost, signed.UploadUrl, &body)
				post.Header.Set("Content-Type", form.FormDataContentType())
				uploaded, err := httpClient.Do(post)
				if err != nil {
					failures.Add(1)
					continue
				}
				_, _ = io.Copy(io.Discard, uploaded.Body)
				uploaded.Body.Close()
				if uploaded.StatusCode >= 300 {
					failures.Add(1)
					continue
				}
				data, _ := json.Marshal(common.OssObjectReq{ObjectKey: signed.ObjectKey})
				post, _ = http.NewRequest(http.MethodPost, target+"/files/complete", bytes.NewReader(data))
				post.Header.Set("Authorization", tokens[tenant].AccessToken)
				post.Header.Set("Content-Type", "application/json")
				post.Header.Set("X-Forwarded-For", fmt.Sprintf("12.%d.%d.1", w+1, i+1))
				completed, err := httpClient.Do(post)
				if err != nil {
					failures.Add(1)
					continue
				}
				_, _ = io.Copy(io.Discard, completed.Body)
				completed.Body.Close()
				if completed.StatusCode != 200 {
					failures.Add(1)
				}
				_ = client.RemoveObject(context.Background(), global.Config.Oss.BucketName, signed.ObjectKey, minio.RemoveObjectOptions{})
			}
		}(w)
	}
	wg.Wait()
	t.Logf("S3 uploads=128 size=4096 concurrency=16 duration=%s throughput=%.1f/s errors=%d", time.Since(started), 128/time.Since(started).Seconds(), failures.Load())
	if failures.Load() != 0 {
		t.Fatal("upload workload failed")
	}
	// Tenant budget is shared across processes, and another tenant stays usable.
	for i := 0; i < 240; i++ {
		request, _ := http.NewRequest(http.MethodGet, targets[i%2]+"/user/profile", nil)
		request.Header.Set("Authorization", tokens[0].AccessToken)
		request.Header.Set("X-Forwarded-For", fmt.Sprintf("13.1.%d.%d", i/200+1, i%200+1))
		response, err := httpClient.Do(request)
		must(t, err)
		response.Body.Close()
	}
	request, _ := http.NewRequest(http.MethodGet, targets[0]+"/user/profile", nil)
	request.Header.Set("Authorization", tokens[9].AccessToken)
	request.Header.Set("X-Forwarded-For", "14.1.1.1")
	response, err := httpClient.Do(request)
	must(t, err)
	response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatal("neighbor tenant unavailable")
	}
}

type recoveryFixture struct {
	Users   []system.SysUser `json:"users"`
	Tokens  []common.Token   `json:"tokens"`
	Hashes  []string         `json:"hashes"`
	Keys    []string         `json:"keys"`
	Digests []string         `json:"digests"`
	JobID   uint             `json:"job_id"`
}

func TestPreproductionRecoverySeed(t *testing.T) {
	if os.Getenv("PREPROD_RECOVERY") != "seed" {
		t.Skip("explicit recovery seed mode required")
	}
	client := storageClient(t)
	fixture := recoveryFixture{}
	for i := 0; i < 2; i++ {
		tid := util.NextID()
		must(t, bypass().Create(&system.SysTenant{TenantId: tid, Name: "recovery", Code: tid, Status: true, MaxStorageBytes: 1048576}).Error)
		user := system.SysUser{UserId: util.NextID(), Username: "recovery", TenantId: tid, Password: "RecoveryOnly!123", AuthVersion: 7}
		must(t, bypass().Create(&user).Error)
		must(t, bypass().Create(&system.SysUserTenant{UserId: user.UserId, TenantId: tid, Status: true}).Error)
		token, err := (&util.Token{UserId: user.UserId, TenantId: tid, ExpireTime: time.Hour, RefreshExpireTime: 2 * time.Hour}).GenerateToken()
		must(t, err)
		payload := []byte(fmt.Sprintf("recovery tenant %d, immutable payload", i))
		signed, err := (&svc.OssService{}).GetPresignedUrl(tenantDB(tid), user.UserId, tid, common.OssReq{FileName: "restore.txt", FileSize: int64(len(payload))})
		must(t, err)
		_, err = client.PutObject(context.Background(), global.Config.Oss.BucketName, signed.ObjectKey, bytes.NewReader(payload), int64(len(payload)), minio.PutObjectOptions{ContentType: "text/plain"})
		must(t, err)
		_, err = storage.Complete(tenantDB(tid), user.UserId, signed.ObjectKey)
		must(t, err)
		fixture.Hashes = append(fixture.Hashes, user.Password)
		fixture.Users = append(fixture.Users, user)
		fixture.Tokens = append(fixture.Tokens, *token)
		fixture.Keys = append(fixture.Keys, signed.ObjectKey)
		fixture.Digests = append(fixture.Digests, fmt.Sprintf("%x", sha256.Sum256(payload)))
		if i == 0 {
			job := &system.SysJob{Name: "restored", JobGroup: "recovery", CronExpression: "0 0 0 * * *", Command: "clean_empty_username_operation_logs", Status: 2, TenantId: tid}
			must(t, bypass().Create(job).Error)
			fixture.JobID = job.Id
			_, err = (&shared.JobScheduler{}).SubmitJob(job)
			must(t, err)
		}
	}
	must(t, client.SetBucketPolicy(context.Background(), global.Config.Oss.BucketName, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::tenant-lifecycle-test/public/avatars/*"]}]}`))
	data, err := json.Marshal(fixture)
	must(t, err)
	must(t, os.WriteFile(os.Getenv("PREPROD_FIXTURE"), data, 0600))
}
func TestPreproductionRecoveryVerify(t *testing.T) {
	if os.Getenv("PREPROD_RECOVERY") != "verify" {
		t.Skip("explicit recovery verify mode required")
	}
	client := storageClient(t)
	data, err := os.ReadFile(os.Getenv("PREPROD_FIXTURE"))
	must(t, err)
	var fixture recoveryFixture
	must(t, json.Unmarshal(data, &fixture))
	for i, user := range fixture.Users {
		var restored system.SysUser
		must(t, tenantDB(user.TenantId).First(&restored, "id = ?", user.UserId).Error)
		if restored.AuthVersion != 7 || restored.Password != fixture.Hashes[i] {
			t.Fatal("credential state differs")
		}
		info, err := util.ValidateAccessToken(context.Background(), fixture.Tokens[i].AccessToken)
		must(t, err)
		if info.UserId != user.UserId || info.TenantId != user.TenantId {
			t.Fatal("restored Redis session differs")
		}
		object, err := client.GetObject(context.Background(), global.Config.Oss.BucketName, fixture.Keys[i], minio.GetObjectOptions{})
		must(t, err)
		payload, err := io.ReadAll(object)
		must(t, err)
		must(t, object.Close())
		if fmt.Sprintf("%x", sha256.Sum256(payload)) != fixture.Digests[i] {
			t.Fatal("object bytes differ")
		}
		charged, limit, err := storage.Usage(tenantDB(user.TenantId), user.TenantId)
		must(t, err)
		if charged != int64(len(payload)) || limit != 1048576 {
			t.Fatal("storage counter differs")
		}
		if _, err := (&svc.OssService{}).GetDownloadUrl(tenantDB(fixture.Users[1-i].TenantId), fixture.Keys[i]); err == nil {
			t.Fatal("restored tenant isolation failed")
		}
		response, err := http.Get("http://" + os.Getenv("OSS_TEST_ENDPOINT") + "/" + global.Config.Oss.BucketName + "/" + fixture.Keys[i])
		must(t, err)
		response.Body.Close()
		if response.StatusCode != 403 {
			t.Fatal("restored private bucket policy differs")
		}
	}
	var run system.SysJobRun
	must(t, bypass().First(&run, "job_id = ?", fixture.JobID).Error)
	if run.Status != "queued" {
		t.Fatal("durable queue differs")
	}
	js := shared.GetJobScheduler()
	defer js.Stop()
	resumed := waitRun(t, fixture.JobID, "succeeded")
	if resumed.Attempts != 1 {
		t.Fatal("restored occurrence ran more than once")
	}
	t.Log("Restored relational ownership, membership, password hash/version, Redis sessions, storage quota, object bytes, private policy, durable job queue")
}

func TestPreproductionJobLoad(t *testing.T) {
	if os.Getenv("PREPROD_RUN") != "1" {
		t.Skip("set PREPROD_RUN=1 and PREPROD_URLS")
	}
	targets := strings.Split(os.Getenv("PREPROD_URLS"), ",")
	if len(targets) != 2 {
		t.Fatal("two backend processes required")
	}
	var ids []uint
	var tokens []string
	for tenant := 0; tenant < 10; tenant++ {
		u := principal(t)
		role := system.SysRole{Name: "job-load", Code: "job-load", TenantId: u.TenantId}
		must(t, bypass().Create(&role).Error)
		var menus []system.SysMenu
		must(t, bypass().Where("scope <> ?", global.PlatformTenantID).Find(&menus).Error)
		must(t, bypass().Model(&role).Association("Menus").Append(menus))
		must(t, bypass().Model(u).Association("Roles").Append(&role))
		token, err := (&util.Token{UserId: u.UserId, TenantId: u.TenantId, ExpireTime: time.Hour, RefreshExpireTime: time.Hour}).GenerateToken()
		must(t, err)
		for j := 0; j < 10; j++ {
			job := system.SysJob{Name: "job-load", JobGroup: "load", CronExpression: "0 0 0 * * *", Command: "clean_empty_username_operation_logs", Status: 2, Concurrent: true, TenantId: u.TenantId}
			must(t, bypass().Create(&job).Error)
			ids = append(ids, job.Id)
			tokens = append(tokens, token.AccessToken)
		}
	}
	client := &http.Client{Timeout: 10 * time.Second}
	var wg sync.WaitGroup
	var failures atomic.Int64
	started := time.Now()
	for worker := 0; worker < 32; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for n := worker; n < 400; n += 32 {
				job := n % 100
				request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/system/jobs/%d/execute", targets[n%2], ids[job]), nil)
				request.Header.Set("Authorization", tokens[job])
				request.Header.Set("X-Forwarded-For", fmt.Sprintf("15.%d.%d.1", n/200+1, n%200+1))
				response, err := client.Do(request)
				if err != nil {
					failures.Add(1)
					continue
				}
				response.Body.Close()
				if response.StatusCode != 204 {
					if failures.Add(1) < 5 {
						t.Logf("job admission status %d", response.StatusCode)
					}
				}
			}
		}(worker)
	}
	wg.Wait()
	if failures.Load() != 0 {
		t.Fatal("job admission failed")
	}
	deadline := time.Now().Add(40 * time.Second)
	for time.Now().Before(deadline) {
		var pending int64
		must(t, bypass().Model(&system.SysJobRun{}).Where("job_id IN ? AND status IN ('queued','retrying')", ids).Count(&pending).Error)
		if pending == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	var runs []system.SysJobRun
	must(t, bypass().Where("job_id IN ?", ids).Find(&runs).Error)
	if len(runs) != 400 {
		t.Fatalf("admitted runs %d", len(runs))
	}
	for _, run := range runs {
		if run.Status != "succeeded" || run.Attempts != 1 {
			t.Fatalf("run not committed exactly once: %+v", run)
		}
	}
	var successTimes int64
	must(t, bypass().Model(&system.SysJob{}).Where("id IN ?", ids).Select("SUM(times)").Scan(&successTimes).Error)
	if successTimes != 400 {
		t.Fatalf("success count %d", successTimes)
	}
	t.Logf("Jobs tenants=10 jobs=100 occurrences=400 replicas=2 request_workers=32 duration=%s committed_throughput=%.1f/s duplicate_effects=0 errors=0", time.Since(started), 400/time.Since(started).Seconds())
}
