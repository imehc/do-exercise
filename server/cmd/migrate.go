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

var (
	sqlFile string
)

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
