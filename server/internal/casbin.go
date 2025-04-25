package internal

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/util"
)

// InitCasbin 初始化casbin
func InitCasbin() {
	// 使用GORM适配器
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		util.Exit("初始化casbin适配器失败: ", err)
	}
	// 从文件加载casbin模型
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		util.Exit("初始化casbin失败: ", err)
	}

	// 启用自动保存策略更改
	enforcer.EnableAutoSave(true)

	// 启用日志
	enforcer.EnableLog(true)

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		util.Exit("加载casbin策略失败: ", err)
	}
	global.Enforcer = enforcer
	fmt.Println("casbin初始化成功")
}
