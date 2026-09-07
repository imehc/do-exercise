package internal

import (
	"reflect"

	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
)

// TenantPlugin 行级租户隔离插件。
// 数据访问层唯一隔离入口：为所有租户隔离表的 SQL 自动追加 tenant_id 条件，
// 并在插入时自动回填 tenant_id，业务代码无需感知租户过滤逻辑。
type TenantPlugin struct{}

func (p *TenantPlugin) Name() string {
	return "tenant:row-scope"
}

func (p *TenantPlugin) Initialize(db *gorm.DB) error {
	db.Callback().Query().Before("gorm:query").Register("tenant:scope", p.scope)
	db.Callback().Update().Before("gorm:update").Register("tenant:scope", p.scope)
	db.Callback().Delete().Before("gorm:delete").Register("tenant:scope", p.scope)
	db.Callback().Create().Before("gorm:create").Register("tenant:fill", p.fill)
	return nil
}

// resolveTenant 解析当前语句归属的租户。
// ok=false 表示该语句无需租户隔离（平台层跨租户操作 / 公共端点无租户上下文）。
func (p *TenantPlugin) resolveTenant(db *gorm.DB) (string, bool) {
	stmt := db.Statement
	if stmt.Context == nil {
		return "", false
	}
	if v := stmt.Context.Value(global.ContextTenantBypassKey); v != nil {
		if bypass, ok := v.(bool); ok && bypass {
			return "", false
		}
	}
	tid, _ := stmt.Context.Value(global.ContextTenantIDKey).(string)
	return global.ResolveTenantID(tid)
}

// tableName 解析语句对应的表名。
// 优先用 Statement.Table；为空时从 Model/Dest 推断（GORM 查询回调执行前
// Statement.Table 可能尚未由 BuildQuerySQL 解析出来）。
func tableName(stmt *gorm.Statement) string {
	if stmt.Table != "" {
		return stmt.Table
	}
	model := stmt.Model
	if model == nil {
		model = stmt.Dest
	}
	return modelTableName(model)
}

// modelTableName 从模型值（含指针/切片包装）解析 TableName
func modelTableName(model interface{}) string {
	if model == nil {
		return ""
	}
	if tn, ok := model.(interface{ TableName() string }); ok {
		return tn.TableName()
	}
	v := reflect.ValueOf(model)
	for v.IsValid() {
		switch v.Kind() {
		case reflect.Ptr:
			if v.IsNil() {
				return ""
			}
			v = v.Elem()
		case reflect.Slice, reflect.Array:
			if v.Type().Elem().Kind() != reflect.Struct {
				return ""
			}
			v = reflect.New(v.Type().Elem()).Elem()
		case reflect.Struct:
			if tn, ok := v.Interface().(interface{ TableName() string }); ok {
				return tn.TableName()
			}
			return ""
		default:
			return ""
		}
	}
	return ""
}

// scope 为查询/更新/删除语句追加租户过滤
func (p *TenantPlugin) scope(db *gorm.DB) {
	if !global.IsTenantScopedTable(tableName(db.Statement)) {
		return
	}
	tid, ok := p.resolveTenant(db)
	if !ok {
		return
	}
	db.Where("tenant_id = ?", tid)
}

// fill 为创建语句回填 tenant_id
func (p *TenantPlugin) fill(db *gorm.DB) {
	if !global.IsTenantScopedTable(tableName(db.Statement)) {
		return
	}
	tid, ok := p.resolveTenant(db)
	if !ok {
		return
	}
	db.Statement.SetColumn("tenant_id", tid)
}
