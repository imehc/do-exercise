package internal

import (
	"fmt"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// InitGorm 初始化PostgreSQL数据库连接
func InitGorm() {
	// 构建DSN连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		global.Config.Database.Host,
		global.Config.Database.Username,
		global.Config.Database.Password,
		global.Config.Database.Database,
		global.Config.Database.Port,
	)

	// 配置GORM
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %v", err))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("获取数据库实例失败: %v", err))
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(global.Config.Database.Pool.MinConnections)
	sqlDB.SetMaxOpenConns(global.Config.Database.Pool.MaxConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(global.Config.Database.Pool.MaxIdleTime) * time.Second)

	// 将数据库实例设置为全局变量
	global.DB = db

	fmt.Println("数据库连接初始化成功")
}
