package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/internal"
	"github.com/imehc/do-exercise/server/migration"
	"github.com/imehc/do-exercise/server/util"
	"github.com/spf13/cobra"
)

var (
	sqlFile    string
	statusOnly bool
)

// appliedTimeLayout `migrate --status` 里执行时间的展示格式
const appliedTimeLayout = "2006-01-02 15:04:05"

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "执行数据库迁移",
	Long: `按版本顺序执行尚未执行的数据库迁移。

已执行的版本登记在 schema_migrations 表，重复执行只会跳过；
每个版本在自己的事务里执行，整轮迁移由 advisory lock 串行化。
--status 只列出各版本的执行状态，不做任何写入。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化配置和数据库连接（AutoMigrate 负责建表建列，迁移只做它做不到的部分）
		internal.InitConfig(configFile)
		internal.InitGorm(true)
		db := global.DB
		if db == nil {
			util.Exit("获取数据库连接失败\n", nil)
		}
		ctx := context.Background()

		if statusOnly {
			// 只读视图不需要种子内容，版本 1 不会被执行
			statuses, err := migration.Statuses(ctx, db, migration.All(""))
			if err != nil {
				util.Exit("读取迁移状态失败: %v\n", err)
			}
			printStatuses(statuses)
			return
		}

		seedSQL := readSeedSQL()
		executed, err := migration.Run(ctx, db, migration.All(seedSQL))
		if err != nil {
			// 失败即中断：已执行的版本已经登记，修完数据后重跑会从断点继续
			for _, m := range executed {
				fmt.Printf("已执行 %d_%s\n", m.Version, m.Name)
			}
			util.Exit("迁移失败: %v\n", err)
		}
		if len(executed) == 0 {
			fmt.Println("数据库已是最新版本，无需迁移")
		} else {
			for _, m := range executed {
				fmt.Printf("已执行 %d_%s\n", m.Version, m.Name)
			}
			fmt.Printf("迁移完成，共执行 %d 个版本\n", len(executed))
		}

		if missing := migration.MissingMenuI18nKeys(db); missing != "" {
			fmt.Printf("以下菜单尚未设置国际化键，界面将回落到菜单名称：%s\n", missing)
		}
	},
}

// readSeedSQL 读取种子 SQL 文件内容（仅版本 1 使用）。
func readSeedSQL() string {
	if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
		util.Exit(fmt.Sprintf("SQL文件不存在: %s\n", sqlFile), nil)
	}
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		util.Exit("读取SQL文件失败: %v\n", err)
	}
	fmt.Printf("种子文件: %s\n", filepath.Base(sqlFile))
	return string(content)
}

// printStatuses 输出各版本的执行状态。
func printStatuses(statuses []migration.Status) {
	pending := 0
	for _, s := range statuses {
		if s.AppliedAt == nil {
			pending++
			fmt.Printf("待执行  %4d  %s\n", s.Version, s.Name)
			continue
		}
		fmt.Printf("已执行  %4d  %s  (%s)\n", s.Version, s.Name, s.AppliedAt.Format(appliedTimeLayout))
	}
	fmt.Printf("共 %d 个版本，%d 个待执行\n", len(statuses), pending)
}

func init() {
	// SQL 文件路径只被版本 1（种子播种）消费，存量库不会用到
	migrateCmd.Flags().StringVar(&sqlFile, "sql", "init.sql", "种子SQL文件路径（仅首次初始化时使用）")
	migrateCmd.Flags().BoolVar(&statusOnly, "status", false, "只查看各版本执行状态，不执行迁移")

	// 添加migrate子命令
	rootCmd.AddCommand(migrateCmd)
}
