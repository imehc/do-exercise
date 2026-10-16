package migration

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
)

// advisoryLockKey 迁移串行化用的 PostgreSQL advisory lock 键。
// 取值只要在本库内唯一且稳定即可（"do_mig" 的 ASCII）；改了它等于换了一把锁，
// 新旧版本的迁移进程会互相看不见对方。
const advisoryLockKey int64 = 0x646F5F6D6967

// 等待迁移锁的上限与重试间隔。上限要能覆盖一次完整迁移的耗时（含大表建索引），
// 又不能长到让部署卡死后没人发现。
const (
	lockWaitTimeout   = 5 * time.Minute
	lockRetryInterval = 2 * time.Second
)

// schemaMigrationsDDL 迁移记录表。
// 这张表自己不能走版本化迁移（它就是版本化的载体），因此只能是幂等 DDL。
const schemaMigrationsDDL = `CREATE TABLE IF NOT EXISTS schema_migrations (
	version    BIGINT PRIMARY KEY,
	name       TEXT NOT NULL,
	applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
)`

// Status 单个版本的执行状态。AppliedAt 为 nil 表示尚未执行。
type Status struct {
	Version   uint64
	Name      string
	AppliedAt *time.Time
}

// Run 按版本升序执行尚未登记的迁移，返回本次实际执行的版本。
//
// 已登记的版本直接跳过——判定依据是 schema_migrations 里的记录，而不是
// 「看看这个索引在不在」这类推断。重复执行 Run 是安全的（附录 B-35）。
func Run(ctx context.Context, db *gorm.DB, migrations []Migration) ([]Migration, error) {
	if err := Validate(migrations); err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("数据库连接为空")
	}

	// 迁移天然是跨租户的 DDL / 全表数据修补，显式声明旁路租户插件。
	// 插件在没有 tenantId 的上下文里本来也不加条件，但那是它的默认行为，
	// 不是这里的意图——默认行为哪天变了，迁移不该跟着变。
	db = db.WithContext(context.WithValue(ctx, global.ContextTenantBypassKey, true))

	unlock, err := lock(ctx, db)
	if err != nil {
		return nil, err
	}
	defer unlock()

	if err := db.Exec(schemaMigrationsDDL).Error; err != nil {
		return nil, fmt.Errorf("创建 schema_migrations 表失败: %w", err)
	}
	applied, err := appliedVersions(db)
	if err != nil {
		return nil, err
	}

	executed := make([]Migration, 0, len(migrations))
	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}
		// 变更与登记放同一个事务：要么都成立，要么都不成立，不会出现
		// 「改了库但没记上」（下次重放）或「记上了但没改」（永久跳过）。
		// 例外是序列操作——setval 不受回滚约束，所以那类迁移自身必须幂等。
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := m.Up(tx); err != nil {
				return err
			}
			return tx.Exec(
				"INSERT INTO schema_migrations (version, name) VALUES (?, ?)",
				m.Version, m.Name,
			).Error
		}); err != nil {
			return executed, fmt.Errorf("迁移 %d_%s 失败: %w", m.Version, m.Name, err)
		}
		executed = append(executed, m)
	}
	return executed, nil
}

// Statuses 汇总每个版本的执行状态，供只读查看。
// 记录表不存在时视为全部待执行，且不建表——只读命令不该留下副作用。
func Statuses(ctx context.Context, db *gorm.DB, migrations []Migration) ([]Status, error) {
	if db == nil {
		return nil, errors.New("数据库连接为空")
	}
	db = db.WithContext(context.WithValue(ctx, global.ContextTenantBypassKey, true))

	records := make(map[uint64]time.Time)
	var tableCount int64
	if err := db.Raw(
		"SELECT COUNT(*) FROM pg_class WHERE relname = 'schema_migrations' AND relkind = 'r'",
	).Scan(&tableCount).Error; err != nil {
		return nil, fmt.Errorf("检查 schema_migrations 表失败: %w", err)
	}
	if tableCount > 0 {
		var rows []struct {
			Version   uint64
			AppliedAt time.Time
		}
		if err := db.Raw("SELECT version, applied_at FROM schema_migrations").Scan(&rows).Error; err != nil {
			return nil, fmt.Errorf("读取迁移记录失败: %w", err)
		}
		for _, r := range rows {
			records[r.Version] = r.AppliedAt
		}
	}

	statuses := make([]Status, 0, len(migrations))
	for _, m := range migrations {
		status := Status{Version: m.Version, Name: m.Name}
		if at, ok := records[m.Version]; ok {
			appliedAt := at
			status.AppliedAt = &appliedAt
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

// appliedVersions 读取已登记的版本号集合。
func appliedVersions(db *gorm.DB) (map[uint64]bool, error) {
	var versions []uint64
	if err := db.Raw("SELECT version FROM schema_migrations").Scan(&versions).Error; err != nil {
		return nil, fmt.Errorf("读取迁移记录失败: %w", err)
	}
	applied := make(map[uint64]bool, len(versions))
	for _, v := range versions {
		applied[v] = true
	}
	return applied, nil
}

// lock 获取整轮迁移的排他锁，返回释放函数。
//
// 用 pg_advisory_xact_lock 而不是会话级的 pg_advisory_lock：database/sql 的连接是
// 池化的，把连接还回池里并不结束会话，会话级锁会一直留在那条连接上，直到进程退出——
// 一个没有释放路径的锁。事务级锁在持锁事务结束时必定释放，没有泄漏可能。
// 代价是整轮迁移期间独占一条连接，因此连接池上限必须大于 1（当前配置是 25）。
//
// 拿不到锁时**等待而不是立刻失败**：`server/Dockerfile` 的启动命令是
// `migrate && server`，多副本同时启动时必然有副本抢不到锁，直接报错会让那个容器
// 起不来。等前一个迁移跑完再进来，看到版本都已登记，自然空转通过。
// 但等待有上限——锁被卡住（例如某个迁移进程挂在半路）时必须报错退出，
// 而不是把部署无限期挂住。
func lock(ctx context.Context, db *gorm.DB) (func(), error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}
	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("开启迁移锁事务失败: %w", err)
	}

	deadline := time.Now().Add(lockWaitTimeout)
	for attempt := 0; ; attempt++ {
		var acquired bool
		if err := tx.QueryRowContext(ctx, "SELECT pg_try_advisory_xact_lock($1)", advisoryLockKey).Scan(&acquired); err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("获取迁移锁失败: %w", err)
		}
		if acquired {
			return func() { _ = tx.Rollback() }, nil
		}
		if attempt == 0 {
			fmt.Println("另一个迁移进程正在执行，等待其完成...")
		}
		if time.Now().After(deadline) {
			_ = tx.Rollback()
			return nil, fmt.Errorf("等待迁移锁超过 %s 仍未获得，可能有迁移进程卡住", lockWaitTimeout)
		}
		select {
		case <-ctx.Done():
			_ = tx.Rollback()
			return nil, ctx.Err()
		case <-time.After(lockRetryInterval):
		}
	}
}
