//go:build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/internal"
	"gorm.io/gorm"
)

var db *gorm.DB

// TestMain 搭建隔离回归测试的运行环境：
// go test 在包目录（server/integration）运行，先切到 server/ 使相对路径可解析，
// 然后读测试配置并连接独立测试库 do_exercise_test（AutoMigrate 已注册 TenantPlugin）。
func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		panic(err)
	}

	// 连接信息允许外部覆盖；默认指向独立测试库/测试 Redis db，避免污染开发环境。
	setDefaultEnv("POSTGRES_HOST", "localhost")
	setDefaultEnv("POSTGRES_PORT", "5432")
	setDefaultEnv("POSTGRES_USER", "admin")
	setDefaultEnv("POSTGRES_PASSWORD", "admin2025")
	setDefaultEnv("POSTGRES_DB", "do_exercise_test")
	setDefaultEnv("REDIS_HOST", "localhost")
	setDefaultEnv("REDIS_PORT", "6379")
	setDefaultEnv("REDIS_PASSWORD", "")
	setDefaultEnv("REDIS_DATABASE", "15")

	internal.InitConfig("integration/testdata/config.yaml")
	internal.InitGorm(true) // 连接测试库并 AutoMigrate（内部注册 TenantPlugin）
	db = global.DB

	code := m.Run()
	os.Exit(code)
}

func setDefaultEnv(key, val string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, val)
	}
}

// bypass 返回绕过租户隔离的 DB 会话：用于种子数据，插件不覆盖/不追加 tenant_id。
func bypass() *gorm.DB {
	ctx := context.WithValue(context.Background(), global.ContextTenantBypassKey, true)
	return db.WithContext(ctx)
}

// tenantDB 返回绑定指定租户上下文的 DB 会话：读写均按该租户隔离。
func tenantDB(tenantID string) *gorm.DB {
	ctx := context.WithValue(context.Background(), global.ContextTenantIDKey, tenantID)
	return db.WithContext(ctx)
}

// truncateScopedTables 清空四张租户隔离表并重置自增序列，保证每个用例幂等。
func truncateScopedTables(t *testing.T) {
	t.Helper()
	if err := bypass().Exec(`TRUNCATE sys_user, sys_role, sys_job, sys_operation_log RESTART IDENTITY CASCADE`).Error; err != nil {
		t.Fatalf("truncate scoped tables: %v", err)
	}
}
