package internal

import (
	"time"

	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Gorm = new(_gorm)

type _gorm struct{}

func (g *_gorm) Config() *gorm.Config {
	general := global.CONFIG.Datebase
	return &gorm.Config{
		Logger: logger.New(NewWriter(general), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      general.LogLevel(),
			Colorful:      true,
		}),
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}
