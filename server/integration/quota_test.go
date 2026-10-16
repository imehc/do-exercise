//go:build integration

package integration

import (
	"strings"
	"sync"
	"testing"

	"github.com/imehc/do-exercise/server/internal"
	sysmodel "github.com/imehc/do-exercise/server/model/system"
	sysreq "github.com/imehc/do-exercise/server/model/system/request"
	syssvc "github.com/imehc/do-exercise/server/service/system"
)

var (
	initEnforcerOnce sync.Once
)

// ensureEnforcer 让走完整服务层的用例可用到 Casbin。进程内只初始化一次。
func ensureEnforcer() {
	initEnforcerOnce.Do(func() {
		internal.InitCasbin()
	})
}

// seedQuotaTenant 种一个带配额上限的租户，并登记清理。
func seedQuotaTenant(t *testing.T, tenantID string, maxUsers, maxTasks int) {
	t.Helper()
	tenant := &sysmodel.SysTenant{
		TenantId: tenantID,
		Name:     "quota-" + tenantID,
		Code:     "quota-" + tenantID,
		Status:   true,
		MaxUsers: maxUsers,
		MaxTasks: maxTasks,
	}
	if err := bypass().Create(tenant).Error; err != nil {
		t.Fatalf("seed quota tenant %s: %v", tenantID, err)
	}
	t.Cleanup(func() {
		_ = bypass().Where("tenant_id = ?", tenantID).Delete(&sysmodel.SysTenant{}).Error
	})
}

// TestTenantUserQuota 验证用户数配额：上限内可建、超限被拒、0 表示不限。
func TestTenantUserQuota(t *testing.T) {
	truncateScopedTables(t)
	ensureEnforcer()

	const capped = "quota-user-cap"
	seedQuotaTenant(t, capped, 1, 0)
	svc := &syssvc.SysUserService{}
	create := func(tenantID, name string) error {
		_, err := svc.Create(tenantDB(tenantID), sysreq.CreateSysUserReq{
			Username: name,
			Nickname: name,
			Password: "quotaPassw0rd!",
		})
		return err
	}

	if err := create(capped, "quota-u1"); err != nil {
		t.Fatalf("配额内创建第 1 个用户应成功: %v", err)
	}
	err := create(capped, "quota-u2")
	if err == nil || !strings.Contains(err.Error(), "tenantUserQuotaReached") {
		t.Fatalf("配额满后再建用户应被拒，期望 tenantUserQuotaReached，实际: %v", err)
	}

	// 0 = 不限：对照租户能连续建两个用户
	const open = "quota-user-open"
	seedQuotaTenant(t, open, 0, 0)
	if err := create(open, "quota-o1"); err != nil {
		t.Fatalf("不限配额租户第 1 个用户应成功: %v", err)
	}
	if err := create(open, "quota-o2"); err != nil {
		t.Fatalf("不限配额租户第 2 个用户应成功: %v", err)
	}
}

// TestTenantJobQuota 验证定时任务数配额：超限被拒。
func TestTenantJobQuota(t *testing.T) {
	truncateScopedTables(t)

	const capped = "quota-job-cap"
	seedQuotaTenant(t, capped, 0, 1)
	svc := &syssvc.SysJobService{}
	create := func(tenantID, name string) error {
		_, err := svc.Create(tenantDB(tenantID), sysreq.CreateSysJobReq{
			Name:           name,
			JobGroup:       "default",
			CronExpression: "0 0 * * *",
			Command:        "clean_empty_username_operation_logs",
			Status:         2, // 暂停态：不走调度器
		})
		return err
	}

	if err := create(capped, "quota-job-1"); err != nil {
		t.Fatalf("配额内创建第 1 个任务应成功: %v", err)
	}
	err := create(capped, "quota-job-2")
	if err == nil || !strings.Contains(err.Error(), "tenantJobQuotaReached") {
		t.Fatalf("配额满后再建任务应被拒，期望 tenantJobQuotaReached，实际: %v", err)
	}
}
