package initialize

import (
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/initialize/internal"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres() {
	p := global.CONFIG.Datebase
	if p.DbName == "" {
		global.LOG.Error("pgsql数据库连接失败")
		return
	}
	pgsqlConfig := postgres.Config{
		DSN:                  p.Dsn(),
		PreferSimpleProtocol: false,
	}
	if db, err := gorm.Open(postgres.New(pgsqlConfig), internal.Gorm.Config()); err != nil {
		global.LOG.Error("pgsql数据库连接失败", zap.Any("err", err))
		return
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
		global.DB = db
		return
	}
}
