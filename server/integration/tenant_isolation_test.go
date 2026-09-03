//go:build integration

package integration

import (
	"testing"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"gorm.io/gorm"
)

const (
	tenantA = "ta"
	tenantB = "tb"
)

// expectedScopedTables 是 tenantScopedTables 当前应登记的租户隔离表。
// 新增隔离表时必须同步补到 server/global/constant.go 的 tenantScopedTables 与本表：
// 漏登 map → TestTenantScopedTableRegistry 红；漏进本表 → 新表隔离没有被任何用例覆盖。
var expectedScopedTables = []string{"sys_user", "sys_role", "sys_job", "sys_operation_log"}

// expectedSharedTables 明确不作为租户隔离的共享表（防止被误登记成隔离表）。
var expectedSharedTables = []string{"sys_tenant", "sys_api", "sys_menu", "casbin_rule"}

// TestTenantScopedTableRegistry 锁住「登记表 == 期望集合」：既不许漏登记隔离表，
// 也不许把共享表误登记。这是「漏登记一张表就泄漏」这条回归的第一道防线。
func TestTenantScopedTableRegistry(t *testing.T) {
	for _, name := range expectedScopedTables {
		if !global.IsTenantScopedTable(name) {
			t.Errorf("期望 %s 参与租户隔离，但 tenantScopedTables 未登记", name)
		}
	}
	for _, name := range expectedSharedTables {
		if global.IsTenantScopedTable(name) {
			t.Errorf("期望 %s 为共享表（不做行级隔离），却被登记进 tenantScopedTables", name)
		}
	}
}

// scopedFixture 描述单张隔离表的种子数据与用于区分的标记列。
type scopedFixture struct {
	table     string
	markerCol string
	seedA     map[string]any // 租户 A 的一行（绕过插件写入，显式 tenant_id）
	seedB     map[string]any // 租户 B 的一行
}

// TestTenantIsolationPerTable 对每张隔离表验证：租户 A/B 各自只见本租户的行，
// 且跨租户读取被插件挡住（漏掉 tenantScopedTables 登记就会出现 B 的行在 A 里可见）。
func TestTenantIsolationPerTable(t *testing.T) {
	fixtures := []scopedFixture{
		{
			table: "sys_user", markerCol: "username",
			seedA: map[string]any{"id": "isoA-user", "username": "isoA", "password": "x", "tenant_id": tenantA},
			seedB: map[string]any{"id": "isoB-user", "username": "isoB", "password": "x", "tenant_id": tenantB},
		},
		{
			table: "sys_role", markerCol: "code",
			seedA: map[string]any{"name": "isoA-role", "code": "isoA", "tenant_id": tenantA},
			seedB: map[string]any{"name": "isoB-role", "code": "isoB", "tenant_id": tenantB},
		},
		{
			table: "sys_job", markerCol: "name",
			seedA: map[string]any{"name": "isoA", "job_group": "g", "cron_expression": "0 0 * * *", "command": "echo a", "tenant_id": tenantA},
			seedB: map[string]any{"name": "isoB", "job_group": "g", "cron_expression": "0 0 * * *", "command": "echo b", "tenant_id": tenantB},
		},
		{
			table: "sys_operation_log", markerCol: "username",
			seedA: map[string]any{"username": "isoA", "tenant_id": tenantA},
			seedB: map[string]any{"username": "isoB", "tenant_id": tenantB},
		},
	}

	for _, fx := range fixtures {
		fx := fx
		t.Run(fx.table, func(t *testing.T) {
			truncateScopedTables(t)
			if err := bypass().Table(fx.table).Create(fx.seedA).Error; err != nil {
				t.Fatalf("seed A into %s: %v", fx.table, err)
			}
			if err := bypass().Table(fx.table).Create(fx.seedB).Error; err != nil {
				t.Fatalf("seed B into %s: %v", fx.table, err)
			}

			// 各租户只见本租户行
			if n := tenantRows(t, tenantDB(tenantA), fx.table, fx.markerCol, "isoA"); n != 1 {
				t.Errorf("%s: 租户A应见自己的 isoA 行，实际 %d", fx.table, n)
			}
			if n := tenantRows(t, tenantDB(tenantB), fx.table, fx.markerCol, "isoB"); n != 1 {
				t.Errorf("%s: 租户B应见自己的 isoB 行，实际 %d", fx.table, n)
			}

			// 跨租户读取必须为空（漏登记/漏过滤时这里会泄漏出对方行）
			if n := tenantRows(t, tenantDB(tenantA), fx.table, fx.markerCol, "isoB"); n != 0 {
				t.Errorf("%s: 租户A不应看到租户B的 isoB 行，实际 %d（跨租户泄漏）", fx.table, n)
			}
			if n := tenantRows(t, tenantDB(tenantB), fx.table, fx.markerCol, "isoA"); n != 0 {
				t.Errorf("%s: 租户B不应看到租户A的 isoA 行，实际 %d（跨租户泄漏）", fx.table, n)
			}
		})
	}
}

// tenantRows 统计绑定租户的会话里命中标记值的行数（插件会追加 tenant_id 过滤）。
func tenantRows(t *testing.T, tx *gorm.DB, table, markerCol, val string) int {
	t.Helper()
	var rows []map[string]any
	if err := tx.Table(table).Where(markerCol+" = ?", val).Find(&rows).Error; err != nil {
		t.Fatalf("query %s.%s=%s: %v", table, markerCol, val, err)
	}
	return len(rows)
}

// TestTenantFillOnCreate 验证创建回填：在租户上下文中创建行，tenant_id 被插件回填；
// 用于确认新增隔离表在登记后，写入侧也会自动带上租户维度。
func TestTenantFillOnCreate(t *testing.T) {
	truncateScopedTables(t)

	job := &system.SysJob{
		Name:           "fill-job",
		JobGroup:       "g",
		CronExpression: "0 0 * * *",
		Command:        "echo hi",
	}
	if err := tenantDB(tenantA).Create(job).Error; err != nil {
		t.Fatalf("create sys_job under tenant A: %v", err)
	}

	var got string
	if err := bypass().Table("sys_job").Select("tenant_id").Where("name = ?", job.Name).Scan(&got).Error; err != nil {
		t.Fatalf("read back sys_job tenant_id: %v", err)
	}
	if got != tenantA {
		t.Errorf("sys_job 创建时应回填 tenant_id=%s，实际 %q", tenantA, got)
	}
}
