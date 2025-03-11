package initialize

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/imehc/do-exercise/server/global"
)

// 首次填充数据到数据库
func InitData() {
	// 读取SQL文件
	sqlFile := filepath.Join("config", "db.sql")
	sqlContent, err := os.ReadFile(sqlFile)
	if err != nil {
		panic(fmt.Sprintf("Error reading SQL file: %s", err.Error()))
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 执行SQL文件
	if err := tx.Exec(string(sqlContent)).Error; err != nil {
		tx.Rollback()
		panic(fmt.Sprintf("Error executing SQL file: %s", err.Error()))
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		panic(fmt.Sprintf("Error committing transaction: %s", err.Error()))
	}

	fmt.Println("Execution of SQL file succeeded!!")
}
