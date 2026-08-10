package internal

import (
	"fmt"
	"time"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// InitGorm 初始化PostgreSQL数据库连接
func InitGorm(isAutoMigrate bool) {
	// 构建DSN连接字符串；sslmode 由配置驱动，避免默认明文
	sslMode := global.Config.Database.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		global.Config.Database.Host,
		global.Config.Database.Username,
		global.Config.Database.Password,
		global.Config.Database.Database,
		global.Config.Database.Port,
		sslMode,
	)
	if connectTimeout := global.Config.Database.Pool.ConnectionTimeout; connectTimeout > 0 {
		dsn += fmt.Sprintf(" connect_timeout=%d", connectTimeout)
	}

	// 配置GORM
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		util.Exit("连接数据库失败: ", err)
	}

	// 自动迁移数据库表
	if isAutoMigrate {
		err = db.AutoMigrate(
			system.SysUser{},
			system.SysMenu{},
			system.SysRole{},
			system.SysApi{},
			system.SysOperationLog{},
			system.SysJob{},
			gormadapter.CasbinRule{},
		)
		if err != nil {
			util.Exit("自动迁移数据库表失败: ", err)
		}
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		util.Exit("获取数据库实例失败: ", err)
	}

	// 设置连接池参数。
	// MaxIdleTime 喂给 SetConnMaxIdleTime，MaxConnLifetime 单独配置——
	// 此前把 MaxIdleTime 误当成连接存活时间，导致每条连接每 5 分钟被硬杀一次。
	sqlDB.SetMaxIdleConns(global.Config.Database.Pool.MinConnections)
	sqlDB.SetMaxOpenConns(global.Config.Database.Pool.MaxConnections)
	sqlDB.SetConnMaxIdleTime(time.Duration(global.Config.Database.Pool.MaxIdleTime) * time.Second)

	maxLifetime := global.Config.Database.Pool.MaxConnLifetime
	if maxLifetime <= 0 {
		maxLifetime = 3600 // 默认 1 小时
	}
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	// 将数据库实例设置为全局变量
	global.DB = db

	fmt.Println("数据库连接初始化成功")
}
