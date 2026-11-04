//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	commonapi "github.com/imehc/do-exercise/server/api/v1/common"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/middleware"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	req "github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/router"
	authsvc "github.com/imehc/do-exercise/server/service/common"
	svc "github.com/imehc/do-exercise/server/service/system"
	"github.com/imehc/do-exercise/server/util"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func principal(t *testing.T) *system.SysUser {
	t.Helper()
	tid := util.NextID()
	seedQuotaTenant(t, tid, 0, 0)
	user := &system.SysUser{UserId: util.NextID(), Username: "audituser", Password: "AuditPass!123", TenantId: tid, Email: "audit@example.com"}
	must(t, bypass().Create(user).Error)
	must(t, bypass().Create(&system.SysUserTenant{UserId: user.UserId, TenantId: tid, Status: true}).Error)
	t.Cleanup(func() {
		_ = util.RevokeAllUserTokens(user.UserId)
		_ = bypass().Where("tenant_id = ?", tid).Delete(&system.SysUser{}).Error
	})
	return user
}
func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
func issue(t *testing.T, u *system.SysUser) *common.Token {
	t.Helper()
	token, err := (&util.Token{UserId: u.UserId, TenantId: u.TenantId, ExpireTime: time.Minute, RefreshExpireTime: time.Hour}).GenerateToken()
	must(t, err)
	return token
}
func TestSessionIsolationAndRevocation(t *testing.T) {
	a, b := principal(t), principal(t)
	ta, tb := issue(t, a), issue(t, b)
	info, err := util.ValidateAccessToken(context.Background(), tb.AccessToken)
	must(t, err)
	service := svc.SysTokenService{}
	scope := svc.TokenScope{TenantId: a.TenantId}
	rows, err := service.FindAll(scope)
	must(t, err)
	if len(rows) != 1 || rows[0].UserId != a.UserId {
		t.Fatalf("foreign sessions: %+v", rows)
	}
	data, err := json.Marshal(rows)
	must(t, err)
	if strings.Contains(string(data), ta.AccessToken) || strings.Contains(string(data), ta.RefreshToken) || strings.Contains(string(data), `"access_token"`) {
		t.Fatal("credentials exposed")
	}
	if err := service.Delete(req.SysTokenDeleteReq{SessionId: info.SessionId}, scope); err == nil {
		t.Fatal("cross-tenant revoke succeeded")
	}
	if err := service.ModityStatus(req.SysTokenModityStatusReq{SessionId: info.SessionId, Disabled: boolp(true)}, scope); err == nil {
		t.Fatal("cross-tenant disable succeeded")
	}
	own, err := util.ValidateAccessToken(context.Background(), ta.AccessToken)
	must(t, err)
	must(t, service.ModityStatus(req.SysTokenModityStatusReq{SessionId: own.SessionId, Disabled: boolp(true)}, scope))
	if _, err := util.ValidateAccessToken(context.Background(), ta.AccessToken); err == nil {
		t.Fatal("disabled session accepted")
	}
	if _, err := (&util.Token{}).RefreshToken(ta.RefreshToken); err == nil {
		t.Fatal("disabled session refreshed")
	}
	must(t, service.ModityStatus(req.SysTokenModityStatusReq{SessionId: own.SessionId, Disabled: boolp(false)}, scope))
	must(t, service.Delete(req.SysTokenDeleteReq{SessionId: own.SessionId}, scope))
	if _, err := util.ValidateAccessToken(context.Background(), ta.AccessToken); err == nil {
		t.Fatal("revoked access accepted")
	}
	if _, err := (&util.Token{}).RefreshToken(ta.RefreshToken); err == nil {
		t.Fatal("revoked refresh accepted")
	}
}
func boolp(v bool) *bool { return &v }
func TestRefreshSingleConsumer(t *testing.T) {
	u := principal(t)
	token := issue(t, u)
	var success atomic.Int32
	var winner *common.Token
	var mu sync.Mutex
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			next, err := (&util.Token{ExpireTime: time.Minute}).RefreshToken(token.RefreshToken)
			if err == nil {
				success.Add(1)
				mu.Lock()
				winner = next
				mu.Unlock()
			}
		}()
	}
	close(start)
	wg.Wait()
	if success.Load() != 1 {
		t.Fatalf("refresh consumed %d times", success.Load())
	}
	if _, err := util.ValidateAccessToken(context.Background(), token.AccessToken); err == nil {
		t.Fatal("old access still accepted")
	}
	if _, err := util.ValidateAccessToken(context.Background(), winner.AccessToken); err != nil {
		t.Fatalf("concurrent loser revoked winner: %v", err)
	}
	raw, err := global.Redis.Get(context.Background(), util.PrefixRefreshToken+token.RefreshToken).Bytes()
	must(t, err)
	var rotated model.RefreshTokenInfo
	must(t, json.Unmarshal(raw, &rotated))
	rotated.RotatedAt = time.Now().Add(-time.Minute)
	data, err := json.Marshal(rotated)
	must(t, err)
	must(t, global.Redis.SetArgs(context.Background(), util.PrefixRefreshToken+token.RefreshToken, data, redis.SetArgs{KeepTTL: true}).Err())
	// A later replay is rejected and revokes that session only.
	if _, err := (&util.Token{}).RefreshToken(token.RefreshToken); err == nil {
		t.Fatal("rotated refresh replay accepted")
	}
	if _, err := util.ValidateAccessToken(context.Background(), winner.AccessToken); err == nil {
		t.Fatal("replay failed to revoke compromised session")
	}
}
func TestLivePrincipalInvalidation(t *testing.T) {
	for _, kind := range []string{"membership", "tenant", "expired", "deleted", "password"} {
		t.Run(kind, func(t *testing.T) {
			u := principal(t)
			token := issue(t, u)
			switch kind {
			case "membership":
				must(t, bypass().Model(&system.SysUserTenant{}).Where("user_id = ?", u.UserId).UpdateColumn("status", false).Error)
			case "tenant":
				must(t, bypass().Model(&system.SysTenant{}).Where("tenant_id = ?", u.TenantId).UpdateColumn("status", false).Error)
			case "expired":
				must(t, bypass().Model(&system.SysTenant{}).Where("tenant_id = ?", u.TenantId).UpdateColumn("expire_time", time.Now().Add(-time.Minute)).Error)
			case "deleted":
				must(t, bypass().Where("id = ?", u.UserId).Delete(&system.SysUser{}).Error)
			case "password":
				must(t, bypass().Model(&system.SysUser{}).Where("id = ?", u.UserId).UpdateColumn("auth_version", gorm.Expr("auth_version + 1")).Error)
			}
			if _, err := util.ValidateAccessToken(context.Background(), token.AccessToken); err == nil {
				t.Fatal("invalid principal access accepted")
			}
			if _, err := (&util.Token{}).RefreshToken(token.RefreshToken); err == nil {
				t.Fatal("invalid principal refreshed")
			}
			if _, err := (&util.Token{UserId: u.UserId, TenantId: u.TenantId, ExpireTime: time.Minute, RefreshExpireTime: time.Hour}).GenerateToken(); err == nil && kind != "password" {
				t.Fatal("invalid principal issued token")
			}
			if kind != "password" {
				if _, err := (&authsvc.AuthService{}).Authenticate(bypass(), u.TenantId, u.Username, "AuditPass!123"); err == nil {
					t.Fatal("invalid principal authenticated")
				}
			}
		})
	}
}
func TestPasswordChangeRetainsOnlyCurrentSession(t *testing.T) {
	u := principal(t)
	a, b := issue(t, u), issue(t, u)
	old := "AuditPass!123"
	must(t, (&svc.SysUserService{}).ResetPassword(tenantDB(u.TenantId), u.UserId, req.UpdateSysUserPasswordReq{Password: "ChangedPass!123"}, &old, a.AccessToken))
	if _, err := util.ValidateAccessToken(context.Background(), a.AccessToken); err != nil {
		t.Fatalf("current session lost: %v", err)
	}
	if _, err := util.ValidateAccessToken(context.Background(), b.AccessToken); err == nil {
		t.Fatal("other session retained")
	}
	if _, err := (&util.Token{}).RefreshToken(a.RefreshToken); err != nil {
		t.Fatalf("retained session cannot refresh: %v", err)
	}
}

func TestPlatformPasswordChangeRetainsCurrentSession(t *testing.T) {
	u := principal(t)
	u.Username = "p" + u.UserId[len(u.UserId)-9:]
	u.TenantId = global.PlatformTenantID
	u.IsSuperAdmin = true
	u.MustChangePassword = true
	must(t, bypass().Model(u).Updates(map[string]any{"username": u.Username, "email": "", "tenant_id": u.TenantId, "is_super_admin": true, "must_change_password": true}).Error)
	issuePlatform := func() *common.Token {
		token, err := (&util.Token{UserId: u.UserId, TenantId: u.TenantId, AuthorizedTenantIds: []string{}, ExpireTime: time.Hour, RefreshExpireTime: 2 * time.Hour}).GenerateToken()
		must(t, err)
		return token
	}
	a, b := issuePlatform(), issuePlatform()
	info, err := util.ValidateAccessToken(context.Background(), a.AccessToken)
	must(t, err)
	must(t, util.SetSessionDisabled(context.Background(), info.SessionId, true))
	if _, err := util.ValidateAccessToken(context.Background(), a.AccessToken); err == nil {
		t.Fatal("disabled platform session authenticated")
	}
	must(t, util.SetSessionDisabled(context.Background(), info.SessionId, false))
	_, err = util.ValidateAccessToken(context.Background(), a.AccessToken)
	must(t, err)
	old := "AuditPass!123"
	actor := context.WithValue(context.Background(), global.ContextTenantIDKey, u.TenantId)
	actor = context.WithValue(actor, global.ContextIsSuperAdminKey, true)
	actor = context.WithValue(actor, global.ContextUserIDKey, u.UserId)
	must(t, (&svc.SysUserService{}).ResetPassword(db.WithContext(actor), u.UserId, req.UpdateSysUserPasswordReq{Password: "ChangedPass!123"}, &old, a.AccessToken))
	retained, err := util.ValidateAccessToken(context.Background(), a.AccessToken)
	must(t, err)
	if retained.MustChangePassword || !retained.IsSuperAdmin || len(retained.AuthorizedTenantIds) != 0 {
		t.Fatal("platform principal changed after password rotation")
	}
	if _, err := util.ValidateAccessToken(context.Background(), b.AccessToken); err == nil {
		t.Fatal("other platform session retained")
	}
	fresh, err := (&util.Token{}).RefreshToken(a.RefreshToken)
	must(t, err)
	_, err = util.ValidateAccessToken(context.Background(), fresh.AccessToken)
	must(t, err)
}
func TestEmailCandidatesExcludeInvalidTenants(t *testing.T) {
	u := principal(t)
	service := authsvc.UserService{}
	users, err := service.FindUsersByEmail(bypass(), u.Email)
	must(t, err)
	found := false
	for _, item := range users {
		found = found || item.UserId == u.UserId
	}
	if !found {
		t.Fatal("active user missing")
	}
	must(t, bypass().Model(&system.SysTenant{}).Where("tenant_id = ?", u.TenantId).UpdateColumn("expire_time", time.Now().Add(-time.Minute)).Error)
	users, _ = service.FindUsersByEmail(bypass(), u.Email)
	for _, item := range users {
		if item.UserId == u.UserId {
			t.Fatal("expired user listed")
		}
	}
}
func TestPlatformServiceAndRouteGuards(t *testing.T) {
	a, b := principal(t), principal(t)
	userService := svc.SysUserService{}
	if err := userService.AssignUserToTenant(tenantDB(a.TenantId), a.UserId, b.TenantId); err == nil {
		t.Fatal("tenant copied user to another tenant")
	}
	if _, err := userService.ListAssignableTenants(tenantDB(a.TenantId), a.UserId); err == nil {
		t.Fatal("tenant enumerated foreign tenants")
	}
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(func(c *gin.Context) { c.Set("tenantId", a.TenantId); c.Set("isSuperAdmin", true) }, middleware.PlatformMiddleware())
	engine.GET("/platform", func(c *gin.Context) { c.Status(204) })
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/platform", nil))
	if response.Code != http.StatusForbidden {
		t.Fatalf("forged super flag accepted: %d", response.Code)
	}
}
func TestTenantScopeScanAndImmutableOwnership(t *testing.T) {
	a, b := principal(t), principal(t)
	dbA := tenantDB(a.TenantId)
	var users []system.SysUser
	must(t, dbA.Model(&system.SysUser{}).Scan(&users).Error)
	if len(users) != 1 || users[0].UserId != a.UserId {
		t.Fatalf("Scan leaked: %+v", users)
	}
	users = nil
	must(t, dbA.Table("sys_user AS u").Select("u.*").Scan(&users).Error)
	if len(users) != 1 || users[0].UserId != a.UserId {
		t.Fatal("aliased Scan leaked")
	}
	if err := db.Model(&system.SysUser{}).Find(&users).Error; err == nil {
		t.Fatal("missing context read succeeded")
	}
	if err := dbA.Raw("SELECT * FROM sys_user").Scan(&users).Error; err == nil {
		t.Fatal("raw query bypass succeeded")
	}
	if err := dbA.Model(&system.SysUser{}).Where("id = ?", a.UserId).Update("tenant_id", b.TenantId).Error; err == nil {
		t.Fatal("tenant reassignment succeeded")
	}
	must(t, dbA.Model(&system.SysUser{}).Where("id = ?", b.UserId).Update("nickname", "hacked").Error)
	copy := *b
	copy.TenantId = a.TenantId
	copy.Nickname = "hacked"
	copy.Roles = nil
	must(t, dbA.Clauses(clause.OnConflict{UpdateAll: true}).Create(&copy).Error)
	var foreign system.SysUser
	must(t, bypass().First(&foreign, "id = ?", b.UserId).Error)
	if foreign.TenantId != b.TenantId || foreign.Nickname == "hacked" {
		t.Fatal("upsert changed foreign row")
	}
}
func TestConcurrentJobQuota(t *testing.T) {
	tid := util.NextID()
	seedQuotaTenant(t, tid, 0, 1)
	var successes atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := (&svc.SysJobService{}).Create(tenantDB(tid), req.CreateSysJobReq{Name: fmt.Sprintf("quota-%d", i), JobGroup: "g", CronExpression: "0 0 0 * * *", Command: "clean_empty_username_operation_logs", Status: 2})
			if err == nil {
				successes.Add(1)
			}
		}(i)
	}
	wg.Wait()
	if successes.Load() != 1 {
		t.Fatalf("quota admitted %d jobs", successes.Load())
	}
}
func TestJobExecutionLiveStateAndDeduplication(t *testing.T) {
	u := principal(t)
	job := &system.SysJob{Name: "audit", JobGroup: "g", CronExpression: "0 0 0 * * *", Command: "clean_empty_username_operation_logs", Status: 1, TenantId: u.TenantId}
	must(t, bypass().Create(job).Error)
	log := &system.SysOperationLog{TenantId: u.TenantId}
	must(t, bypass().Create(log).Error)
	js := &shared.JobScheduler{}
	slot := time.Now().Truncate(time.Second)
	must(t, js.ExecuteScheduledJob(job, slot))
	must(t, js.ExecuteScheduledJob(job, slot))
	var got system.SysJob
	must(t, bypass().First(&got, job.Id).Error)
	if got.Times != 1 {
		t.Fatalf("same occurrence ran %d times", got.Times)
	}
	must(t, bypass().Model(&system.SysJob{}).Where("id = ?", job.Id).UpdateColumn("status", 2).Error)
	if err := js.ExecuteScheduledJob(job, slot.Add(time.Second)); err == nil {
		t.Fatal("stopped job scheduled")
	}
	must(t, bypass().Model(&system.SysTenant{}).Where("tenant_id = ?", u.TenantId).UpdateColumn("status", false).Error)
	if err := js.ExecuteJob(job); err == nil {
		t.Fatal("disabled tenant job executed")
	}
	var count int64
	must(t, bypass().Model(&system.SysJob{}).Where("id = ?", job.Id).Select("times").Scan(&count).Error)
	if count != 1 {
		t.Fatal("rejected job counted")
	}
}
func TestPermissionRevocationIsCommittedDatabaseState(t *testing.T) {
	u := principal(t)
	role := &system.SysRole{Name: "audit", Code: "audit", TenantId: u.TenantId}
	must(t, bypass().Create(role).Error)
	var menu system.SysMenu
	must(t, bypass().First(&menu, 133).Error)
	must(t, bypass().Model(role).Association("Menus").Append(&menu))
	must(t, bypass().Model(u).Association("Roles").Append(role))
	for i := 0; i < 2; i++ {
		ok, err := util.Authorize(tenantDB(u.TenantId).Session(&gorm.Session{NewDB: true}), u.UserId, u.TenantId, "/system/users/123", "PUT")
		must(t, err)
		if !ok {
			t.Fatal("valid permission rejected")
		}
	}
	must(t, bypass().Transaction(func(tx *gorm.DB) error { return tx.Model(role).Association("Menus").Clear() }))
	for i := 0; i < 2; i++ {
		ok, err := util.Authorize(tenantDB(u.TenantId).Session(&gorm.Session{NewDB: true}), u.UserId, u.TenantId, "/system/users/123", "PUT")
		must(t, err)
		if ok {
			t.Fatal("revoked permission cached")
		}
	}
}

func TestOwnershipCannotBeClearedOrOverriddenOnConflict(t *testing.T) {
	u := principal(t)
	copy := *u
	copy.TenantId = ""
	if err := tenantDB(u.TenantId).Save(&copy).Error; err == nil {
		t.Fatal("Save cleared ownership")
	}
	foreign := principal(t)
	value := &system.SysFile{ObjectKey: "conflict-" + u.UserId, TenantId: u.TenantId, UserId: u.UserId, FileName: "f", ContentType: "text/plain", MaxSize: 100}
	must(t, tenantDB(u.TenantId).Create(value).Error)
	if err := tenantDB(u.TenantId).Clauses(clause.OnConflict{DoUpdates: clause.Assignments(map[string]any{"tenant_id": foreign.TenantId})}).Create(value).Error; err == nil {
		t.Fatal("upsert assignment changed ownership")
	}
}

func TestLastAdministratorAndProtectedRole(t *testing.T) {
	u := principal(t)
	role := &system.SysRole{Name: "tenant admin", Code: svc.TenantAdminRoleCode, TenantId: u.TenantId}
	must(t, bypass().Create(role).Error)
	must(t, bypass().Model(u).Association("Roles").Append(role))
	service := svc.SysUserService{}
	if err := service.Delete(tenantDB(u.TenantId), u.UserId); err == nil {
		t.Fatal("last tenant administrator deleted")
	}
	must(t, service.Update(tenantDB(u.TenantId), u.UserId, req.UpdateSysUserReq{Nickname: "changed", RoleIds: []uint{}}))
	var roles []system.SysRole
	must(t, bypass().Model(u).Association("Roles").Find(&roles))
	if len(roles) != 1 || roles[0].Id != role.Id {
		t.Fatal("tenant administrator role removed by profile edit")
	}
	v := principal(t)
	if err := service.Update(tenantDB(u.TenantId), v.UserId, req.UpdateSysUserReq{RoleIds: []uint{role.Id}}); err == nil {
		t.Fatal("foreign user update succeeded")
	}
	if err := (&svc.SysRoleService{}).Delete(tenantDB(u.TenantId), role.Id); err == nil {
		t.Fatal("protected role deleted")
	}
}

func TestDatabaseInvalidationSurvivesRedisCleanupFailure(t *testing.T) {
	u := principal(t)
	token := issue(t, u)
	oldRedis := global.Redis
	global.Redis = redis.NewClient(&redis.Options{Addr: "127.0.0.1:1", MaxRetries: -1, DialTimeout: 50 * time.Millisecond})
	serviceErr := (&svc.SysUserService{}).ResetPassword(tenantDB(u.TenantId), u.UserId, req.UpdateSysUserPasswordReq{Password: "ChangedPass!123"}, nil, "")
	_ = global.Redis.Close()
	global.Redis = oldRedis
	if serviceErr == nil {
		t.Fatal("cleanup failure was not reported")
	}
	if _, err := util.ValidateAccessToken(context.Background(), token.AccessToken); err == nil {
		t.Fatal("old credential accepted after database password change")
	}
	if _, err := (&util.Token{}).RefreshToken(token.RefreshToken); err == nil {
		t.Fatal("old refresh accepted after database password change")
	}
}

func TestTenantJoinScopesEveryTenantTable(t *testing.T) {
	a, b := principal(t), principal(t)
	for _, u := range []*system.SysUser{a, b} {
		role := &system.SysRole{Name: "r", Code: "r", TenantId: u.TenantId}
		must(t, bypass().Create(role).Error)
		must(t, bypass().Model(u).Association("Roles").Append(role))
	}
	var rows []system.SysUser
	must(t, tenantDB(a.TenantId).Table("sys_user_role r").Select("u.*").Joins("JOIN sys_user u ON u.id = r.sys_user_user_id").Scan(&rows).Error)
	if len(rows) != 1 || rows[0].UserId != a.UserId {
		t.Fatal("tenant join leaked foreign rows")
	}
}

func TestPendingLoginBoundToCredentialVersion(t *testing.T) {
	u := principal(t)
	engine := gin.New()
	engine.Use(middleware.PublicAuthContextMiddleware())
	engine.POST("/select", (&commonapi.AuthApi{}).SelectTenant)
	pending := common.LoginSession{Tenants: []string{u.TenantId}, Principals: map[string]common.LoginPrincipal{u.TenantId: {UserId: u.UserId, AuthVersion: u.AuthVersion, TenantVersion: 1}}}
	data, err := json.Marshal(pending)
	must(t, err)
	id := util.NextID()
	must(t, global.Redis.Set(context.Background(), "loginSession:"+id, data, time.Minute).Err())
	must(t, bypass().Model(&system.SysUser{}).Where("id = ?", u.UserId).UpdateColumn("auth_version", gorm.Expr("auth_version + 1")).Error)
	body := fmt.Sprintf(`{"login_session_id":%q,"tenant_id":%q}`, id, u.TenantId)
	for i := 0; i < 2; i++ {
		request := httptest.NewRequest("POST", "/select", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		engine.ServeHTTP(response, request)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("old login proof accepted: %d", response.Code)
		}
	}
	// Fresh proof is consumed once and can complete authentication.
	pending.Principals[u.TenantId] = common.LoginPrincipal{UserId: u.UserId, AuthVersion: u.AuthVersion + 1, TenantVersion: 1}
	data, err = json.Marshal(pending)
	must(t, err)
	id = util.NextID()
	must(t, global.Redis.Set(context.Background(), "loginSession:"+id, data, time.Minute).Err())
	body = fmt.Sprintf(`{"login_session_id":%q,"tenant_id":%q}`, id, u.TenantId)
	request := httptest.NewRequest("POST", "/select", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("fresh proof rejected: %d %s", response.Code, response.Body.String())
	}
}

func TestLegacyMenuCannotGrantPlatformEndpoints(t *testing.T) {
	a, b := principal(t), principal(t)
	token := issue(t, a)
	var menu system.SysMenu
	must(t, bypass().Preload("Apis").First(&menu, 133).Error)
	original := append([]system.SysApi(nil), menu.Apis...)
	defer func() { _ = bypass().Model(&menu).Association("Apis").Replace(original) }()
	var platformApis []system.SysApi
	must(t, bypass().Where("id IN ?", []uint{129, 130}).Find(&platformApis).Error)
	must(t, bypass().Model(&menu).Association("Apis").Append(platformApis))
	role := &system.SysRole{Name: "legacy", Code: "legacy", TenantId: a.TenantId}
	must(t, bypass().Create(role).Error)
	must(t, bypass().Model(role).Association("Menus").Append(&menu))
	must(t, bypass().Model(a).Association("Roles").Append(role))
	engine := router.New()
	for _, test := range []struct{ method, path, body string }{{"GET", "/system/users/" + a.UserId + "/assignable_tenants", ""}, {"PATCH", "/system/users/" + a.UserId + "/tenant", fmt.Sprintf(`{"tenant_id":%q}`, b.TenantId)}} {
		allowed, err := util.Authorize(tenantDB(a.TenantId), a.UserId, a.TenantId, test.path, test.method)
		must(t, err)
		if !allowed {
			t.Fatal("legacy fixture did not grant API")
		}
		request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
		request.Header.Set("Authorization", token.AccessToken)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		engine.ServeHTTP(response, request)
		if response.Code != 403 {
			t.Fatalf("platform guard failed: %d", response.Code)
		}
	}
	var count int64
	must(t, bypass().Model(&system.SysUser{}).Where("username = ? AND tenant_id = ?", a.Username, b.TenantId).Count(&count).Error)
	if count != 1 {
		t.Fatal("cross-tenant account inserted")
	}
}

func TestAuthenticationDependencyFailurePreservesSession(t *testing.T) {
	u := principal(t)
	token := issue(t, u)
	original := global.DB
	t.Cleanup(func() { global.DB = original })
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	_, dbErr := util.LoadSessionPrincipal(db.WithContext(canceled), u.UserId, u.TenantId)
	if !errors.Is(dbErr, util.ErrAuthDependencyUnavailable) {
		t.Fatalf("database error classified as credentials: %v", dbErr)
	}
	global.DB = db.WithContext(canceled)
	_, err := util.ValidateAccessToken(canceled, token.AccessToken)
	if !errors.Is(err, util.ErrAuthDependencyUnavailable) {
		t.Fatalf("dependency error classified as credentials: %v", err)
	}
	global.DB = original
	_, err = util.ValidateAccessToken(context.Background(), token.AccessToken)
	must(t, err)
}
