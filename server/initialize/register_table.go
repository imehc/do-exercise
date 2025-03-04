package initialize

import (
	"os"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"go.uber.org/zap"
)

func InitRegisterTable() {
	db := global.DB
	log := global.LOG

	err := db.AutoMigrate(
		system.User{},
		system.Dept{},
		system.Dict{},
		system.DictData{},
		system.Post{},
		system.Api{},
		system.Menu{},
		system.Role{},
	)

	if err != nil {
		log.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	log.Info("register table success")
}
