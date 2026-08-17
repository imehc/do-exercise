package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/internal"
	"github.com/imehc/do-exercise/server/util"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var sqlFile string

// defaultAdminHash 是 init.sql 播种的默认管理员口令对应的公开 bcrypt 哈希
// （对应 README.md 记载的 @admin2025）。任何仍使用该哈希的账号都必须强制改密。
const defaultAdminHash = "$2a$10$gI7PJi4gyTc.sG2m5ZgbcO/I0E8nLkW2AHhWFxGMaCogU2H/E3YzC"

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "执行数据库初始化",
	Long:  `执行数据库初始化，需要提供SQL文件路径。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查SQL文件是否存在
		if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
			util.Exit(fmt.Sprintf("SQL文件不存在: %s\n", sqlFile), nil)
		}

		// 读取SQL文件内容
		sqlContent, err := os.ReadFile(sqlFile)
		if err != nil {
			util.Exit("读取SQL文件失败:", err)
		}

		// 初始化配置和数据库连接
		internal.InitConfig(configFile)
		// 获取数据库连接
		internal.InitGorm(true)
		db := global.DB
		if db == nil {
			util.Exit("获取数据库连接失败", nil)
		}

		// 默认管理员口令强制轮换。
		// 对存量数据库（已播种过 init.sql，sys_user 有数据）中的默认管理员打上
		// must_change_password 标记，强制其在下次登录时修改公开的默认口令。
		// 该 UPDATE 幂等：仅当密码仍等于公开的默认哈希时才生效，改密后哈希不再匹配。
		result := db.Exec(
			"UPDATE sys_user SET must_change_password = TRUE WHERE password = ? AND must_change_password = FALSE",
			defaultAdminHash,
		)
		if result.Error != nil {
			util.Exit("标记默认管理员强制改密失败: ", result.Error)
		}
		if rows := result.RowsAffected; rows > 0 {
			fmt.Printf("已将 %d 个仍使用默认口令的账号标记为必须改密\n", rows)
		}

		// 软删除与唯一约束冲突修复：
		// 旧版用普通唯一索引，软删的账号会永久占用用户名/邮箱，导致同名/同邮箱无法重建。
		// 现改为 deleted_at IS NULL 的部分唯一索引（由模型 uniqueIndex 标签声明，AutoMigrate 创建）。
		// 邮箱唯一索引已进一步改为 email <> '' 的条件索引，空邮箱不参与唯一性校验。
		// 存量库里的旧索引不会自动消失，这里幂等清理，避免重建账号时撞上旧约束。
		db.Exec("DROP INDEX IF EXISTS idx_sys_user_username")
		db.Exec("DROP INDEX IF EXISTS idx_sys_user_email")
		db.Exec("DROP INDEX IF EXISTS idx_sys_user_email_tenant")
		// 重建邮箱条件唯一索引：仅 email 非空时参与唯一性校验，空邮箱可在同一租户共存多条。
		db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_user_email_tenant ON sys_user (email, tenant_id) WHERE deleted_at IS NULL AND email <> ''")

		// 检查表是否已有数据
		var count int64
		err = db.Raw("SELECT COUNT(*) FROM sys_user").Count(&count).Error
		if err != nil {
			// 如果表不存在，继续执行初始化
			fmt.Println("表不存在，开始执行初始化...")
		} else if count > 0 {
			fmt.Println("数据库已有数据，跳过初始化")
			return
		}

		// 执行SQL
		fmt.Printf("正在执行SQL文件: %s\n", filepath.Base(sqlFile))
		err = db.Session(&gorm.Session{}).Exec(string(sqlContent)).Error
		if err != nil {
			util.Exit("执行SQL失败: ", err)
		}
		fmt.Println("数据库初始化完成")
	},
}

func init() {
	// 添加必需的SQL文件路径标志
	migrateCmd.Flags().StringVar(&sqlFile, "sql", "", "SQL文件路径（必需）")
	migrateCmd.MarkFlagRequired("sql")

	// 添加migrate子命令
	rootCmd.AddCommand(migrateCmd)
}
