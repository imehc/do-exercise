package initialize

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/imehc/do-exercise/server/global"
	"go.uber.org/zap"
)

func InitCasbin() {
	a, err := gormadapter.NewAdapterByDB(global.DB) // 使用 Gorm 适配器
	if err != nil {
		global.LOG.Error("请检查数据库是否连接成功", zap.Error(err))
		return
	}
	m := model.NewModel()

	casbinModelText := `
	[request_definition]
	r = sub, dom, obj, act

	[policy_definition]
	p = sub, dom, obj, act

	[role_definition]
	g = _, _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj &&  r.act == p.act
	`

	m.LoadModelFromText(casbinModelText) // 加载上面的模型
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		zap.L().Error("加载模型失败!", zap.Error(err))
		return
	}
	// 添加策略（示例）角色-资源-动作-租户
	e.AddPolicy("admin", "/user", "(GET|POST)", "tenant_1")
	global.Enforcer = e
}
