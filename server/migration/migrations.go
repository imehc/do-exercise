package migration

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// defaultAdminHash 是 init.sql 播种的默认管理员口令对应的公开 bcrypt 哈希
// （对应 README.md 记载的 @admin2025）。任何仍使用该哈希的账号都必须强制改密。
const defaultAdminHash = "$2a$10$gI7PJi4gyTc.sG2m5ZgbcO/I0E8nLkW2AHhWFxGMaCogU2H/E3YzC"

// All 返回按版本升序排列的全部迁移。
//
// seedSQL 是 init.sql 的内容，由 migrate 命令读入后传进来（版本 1 需要它）。
// 只读场景（`migrate --status`）可以传空串：版本 1 不会被执行，只用于列出版本号。
//
// 表结构本身仍由 AutoMigrate 负责建列（internal/gorm.go），这里只做 AutoMigrate
// 做不到的部分：种子数据、唯一约束口径变更、序列修复、存量数据回填。
func All(seedSQL string) []Migration {
	return []Migration{
		{
			Version: 1,
			Name:    "seed_initial_data",
			Up: func(db *gorm.DB) error {
				// 存量库（sys_user 已有数据）视为迁移基线：种子不重放，直接登记为已执行。
				// 这条判断必须留在迁移内部而不是调用方——它决定的是「这个版本对这个库
				// 意味着什么」，放到外面就变成了又一处没被记录的隐式分支。
				var count int64
				if err := db.Raw("SELECT COUNT(*) FROM sys_user").Scan(&count).Error; err != nil {
					return fmt.Errorf("统计 sys_user 失败: %w", err)
				}
				if count > 0 {
					fmt.Println("数据库已有数据，视为迁移基线，跳过种子播种")
					return nil
				}
				if seedSQL == "" {
					return errors.New("种子 SQL 为空，无法播种（请检查 --sql 指向的文件）")
				}
				if err := db.Exec(seedSQL).Error; err != nil {
					return fmt.Errorf("执行种子 SQL 失败: %w", err)
				}
				fmt.Println("种子数据播种完成")
				return nil
			},
		},
		{
			Version: 2,
			Name:    "force_default_admin_password_rotation",
			Up: func(db *gorm.DB) error {
				// 默认管理员口令强制轮换：对仍在使用公开默认哈希的账号打上
				// must_change_password，强制其下次登录改密。
				// 幂等——改密后哈希不再匹配，这条 UPDATE 自然命中 0 行。
				// 新库的种子已经把该标记写成 TRUE，所以这条只对存量库起作用。
				result := db.Exec(
					"UPDATE sys_user SET must_change_password = TRUE WHERE password = ? AND must_change_password = FALSE",
					defaultAdminHash,
				)
				if result.Error != nil {
					return fmt.Errorf("标记默认管理员强制改密失败: %w", result.Error)
				}
				if rows := result.RowsAffected; rows > 0 {
					fmt.Printf("已将 %d 个仍使用默认口令的账号标记为必须改密\n", rows)
				}
				return nil
			},
		},
		{
			Version: 3,
			Name:    "sys_user_soft_delete_aware_unique_indexes",
			Up: func(db *gorm.DB) error {
				// 软删除与唯一约束冲突修复：
				// 旧版用普通唯一索引，软删的账号会永久占用用户名/邮箱，同名/同邮箱无法重建。
				// 现口径是 deleted_at IS NULL 的部分唯一索引（模型 uniqueIndex 标签声明，
				// 新库由 AutoMigrate 直接建出），邮箱再加 email <> ''，空邮箱不参与唯一性。
				// 存量库里的旧索引不会自动消失，这里显式清理。
				//
				// 先 DROP CONSTRAINT 再 DROP INDEX：如果旧唯一性是以列约束形式存在，
				// 约束会「拥有」背后的索引，直接 DROP INDEX 会被 PostgreSQL 拒绝。
				return execAll(db, []string{
					"ALTER TABLE sys_user DROP CONSTRAINT IF EXISTS uni_sys_user_username",
					"ALTER TABLE sys_user DROP CONSTRAINT IF EXISTS uni_sys_user_email",
					"DROP INDEX IF EXISTS idx_sys_user_username",
					"DROP INDEX IF EXISTS idx_sys_user_email",
					"DROP INDEX IF EXISTS idx_sys_user_email_tenant",
					"CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_user_email_tenant ON sys_user (email, tenant_id) WHERE deleted_at IS NULL AND email <> ''",
				})
			},
		},
		{
			Version: 4,
			Name:    "sys_role_code_partial_unique_index",
			Up: func(db *gorm.DB) error {
				// 角色编码（P2-3）：唯一性口径是「同租户内、未软删的角色之间唯一」。
				// 旧索引不含 deleted_at 条件，删掉一个角色后重建同 code 会撞唯一约束，
				// 而 checkCodeDuplicate 又看不到软删行，两边错位只能得到一个通用写库错误。
				return execAll(db, []string{
					"ALTER TABLE sys_role DROP CONSTRAINT IF EXISTS uni_sys_role_code",
					"DROP INDEX IF EXISTS idx_sys_role_code",
					"DROP INDEX IF EXISTS idx_sys_role_code_tenant",
					"CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_role_code_tenant ON sys_role (code, tenant_id) WHERE deleted_at IS NULL",
				})
			},
		},
		{
			Version: 5,
			Name:    "sys_menu_permission_partial_unique_index",
			Up: func(db *gorm.DB) error {
				// 菜单权限标识（P2-8）：同理。菜单不分租户，因此保持全局唯一，只加软删条件。
				// 旧口径是列上的 UNIQUE 约束（GORM 的 `unique` 标签），先降级为部分唯一索引。
				return execAll(db, []string{
					"ALTER TABLE sys_menu DROP CONSTRAINT IF EXISTS uni_sys_menu_permission",
					"ALTER TABLE sys_menu DROP CONSTRAINT IF EXISTS sys_menu_permission_key",
					"DROP INDEX IF EXISTS uni_sys_menu_permission",
					"CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_menu_permission ON sys_menu (permission) WHERE deleted_at IS NULL",
				})
			},
		},
		{
			Version: 6,
			Name:    "casbin_rule_sequence_setval",
			Up: func(db *gorm.DB) error {
				// Casbin 策略表的自增序列修复：种子把 p 规则写死成 101-141 的显式 id，
				// 而 casbin_rule.id 是自增主键。旧种子没有重置序列，运行期新增策略从 1 开始
				// 发号，加到第 101 条时撞上种子行主键，表现为「角色授权偶发失败」。
				//
				// setval 不受事务回滚约束，所以这条必须自身幂等：取 MAX(id) 推进，
				// 重复执行只会把序列推到同一个位置。
				if err := db.Exec(
					"SELECT setval(pg_get_serial_sequence('casbin_rule','id'), (SELECT COALESCE(MAX(id), 1) FROM casbin_rule))",
				).Error; err != nil {
					return fmt.Errorf("修复 casbin_rule 自增序列失败: %w", err)
				}
				return nil
			},
		},
		{
			Version: 7,
			Name:    "menu_metadata_backfill",
			Up: func(db *gorm.DB) error {
				// 菜单元数据回填：补上 scope / is_system / i18n_key。
				// 运行期的可见范围判定只认 scope，因此这一步必须先于任何租户请求完成。
				// 对新库（版本 1 刚播完种子）和存量库是同一份口径，只有一处实现。
				return backfillMenuMetadata(db)
			},
		},
		{
			Version: 8,
			Name:    "sys_menu_scope_not_null",
			Up: func(db *gorm.DB) error {
				// 把「scope 不为空」从约定升级成约束。
				//
				// 版本 7 已经把 NULL / '' 回填成 both，但只要列还允许 NULL，
				// 运行期的可见性判定就得一直带着 `scope IS NULL OR ...` 这种兜底——
				// 而那个兜底是 fail-open 的：一条 scope 为 NULL 的平台菜单会被
				// 当成租户可见。约束加上之后兜底才可以安全删掉（P1-3 收尾）。
				//
				// 三条语句都幂等：回填是条件更新，SET DEFAULT / SET NOT NULL 重复执行无副作用。
				return execAll(db, []string{
					"UPDATE sys_menu SET scope = 'both' WHERE scope IS NULL OR scope = ''",
					"ALTER TABLE sys_menu ALTER COLUMN scope SET DEFAULT 'both'",
					"ALTER TABLE sys_menu ALTER COLUMN scope SET NOT NULL",
				})
			},
		},
		{
			Version: 9,
			Name:    "sys_user_tenant_membership",
			Up: func(db *gorm.DB) error {
				// 建 sys_user_tenant 成员关系表，并按现有 sys_user(id, tenant_id) 全量回填
				// （每行在册用户一条成员关系）。这次只新增表、任何旧代码都不读它，因此
				// 可以安全先行落地并随时停用；后续版本才把写入与读取切到这张表。
				//
				// 回填口径：
				//   - 只回填未软删的 sys_user 行。软删的用户 = 已移出该租户，不应残留
				//     成员关系；成员关系表达的是「当前在哪些租户里」这一事实。
				//   - 平台保留租户（tenant_id = platform）同样落一行，让成员关系行数与
				//     在册用户数一一对应，便于事后核对回填完整性。
				//
				// 幂等：CREATE IF NOT EXISTS / 部分索引 IF NOT EXISTS / INSERT ON CONFLICT。
				if err := execAll(db, []string{
					`CREATE TABLE IF NOT EXISTS sys_user_tenant (
						user_id     VARCHAR(32)  NOT NULL,
						tenant_id   VARCHAR(32)  NOT NULL,
						status      BOOLEAN      NOT NULL DEFAULT TRUE,
						created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
						created_by  VARCHAR(32),
						updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
						updated_by  VARCHAR(32),
						deleted_at  TIMESTAMPTZ,
						PRIMARY KEY (user_id, tenant_id)
					)`,
					`CREATE INDEX IF NOT EXISTS idx_sys_user_tenant_tenant
						ON sys_user_tenant (tenant_id) WHERE deleted_at IS NULL`,
				}); err != nil {
					return err
				}

				result := db.Exec(
					`INSERT INTO sys_user_tenant
						(user_id, tenant_id, status, created_at, created_by, updated_at, updated_by)
					 SELECT id, tenant_id, TRUE, created_at, created_by, updated_at, updated_by
					   FROM sys_user
					  WHERE deleted_at IS NULL
					    AND tenant_id IS NOT NULL AND tenant_id <> ''
					  ON CONFLICT (user_id, tenant_id) DO NOTHING`,
				)
				if result.Error != nil {
					return fmt.Errorf("回填 sys_user_tenant 失败: %w", result.Error)
				}
				fmt.Printf("sys_user_tenant 建表并回填完成（写入 %d 行成员关系）\n", result.RowsAffected)
				return nil
			},
		},
	}
}

// execAll 顺序执行一组语句，任一失败即返回。
//
// 与旧实现的区别只有一点：错误不再被丢掉。存量库里有脏数据时（例如同租户下
// 重复的角色编码），CREATE UNIQUE INDEX 会失败，迁移随之中断并保留现场——
// 这正是想要的行为，脚本没有资格替运维决定删哪一行（附录 B-37）。
func execAll(db *gorm.DB, statements []string) error {
	for _, sql := range statements {
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("执行失败 [%s]: %w", sql, err)
		}
	}
	return nil
}
