package internal

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	casbinlog "github.com/casbin/casbin/v3/log"
	"github.com/casbin/casbin/v3/model"
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

	m := model.NewModel()

	casbinModelText := `
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[role_definition]
	g = _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
	`

	m.LoadModelFromText(casbinModelText) // 加载上面的模型
	// 从文件加载casbin模型
	// enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		util.Exit("初始化casbin失败: ", err)
	}

	// 启用自动保存策略更改
	enforcer.EnableAutoSave(true)

	// 启用日志
	if !util.IsRelease {
		enforcer.SetLogger(casbinlog.NewDefaultLogger())
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		util.Exit("加载casbin策略失败: ", err)
	}
	global.Enforcer = enforcer
	fmt.Println("casbin初始化成功")
}
