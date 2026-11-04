package internal

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrTenantRequired = errors.New("tenant context is required")
var ErrTenantImmutable = errors.New("tenant_id cannot be changed")
var sqlWords = regexp.MustCompile(`[A-Za-z_][A-Za-z_0-9]*`)
var sqlJoin = regexp.MustCompile(`(?i)\bJOIN\s+["\x60]?([A-Za-z_][A-Za-z_0-9]*)["\x60]?(?:\s+(?:AS\s+)?([A-Za-z_][A-Za-z_0-9]*))?\s+ON\b`)

func touchesTenantTable(sql string) bool {
	for _, word := range sqlWords.FindAllString(sql, -1) {
		if global.IsTenantScopedTable(strings.ToLower(word)) {
			return true
		}
	}
	return false
}

type TenantPlugin struct{}

func (p *TenantPlugin) Name() string { return "tenant:row-scope" }
func (p *TenantPlugin) Initialize(db *gorm.DB) error {
	if err := db.Callback().Query().Before("gorm:query").Register("tenant:scope", p.scope); err != nil {
		return err
	}
	if err := db.Callback().Row().Before("gorm:row").Register("tenant:scope", p.scope); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:before_update").Register("tenant:scope", p.update); err != nil {
		return err
	}
	if err := db.Callback().Delete().Before("gorm:before_delete").Register("tenant:scope", p.scope); err != nil {
		return err
	}
	if err := db.Callback().Raw().Before("gorm:raw").Register("tenant:raw", p.raw); err != nil {
		return err
	}
	return db.Callback().Create().Before("gorm:before_create").Register("tenant:fill", p.fill)
}

func bypassTenant(db *gorm.DB) bool {
	bypass, _ := db.Statement.Context.Value(global.ContextTenantBypassKey).(bool)
	if bypass {
		return true
	}
	tid, _ := db.Statement.Context.Value(global.ContextTenantIDKey).(string)
	super, _ := db.Statement.Context.Value(global.ContextIsSuperAdminKey).(bool)
	return tid == global.PlatformTenantID && super
}

func (p *TenantPlugin) resolveTenant(db *gorm.DB) (string, bool) {
	if bypassTenant(db) {
		return "", false
	}
	tid, _ := db.Statement.Context.Value(global.ContextTenantIDKey).(string)
	if tid == "" || tid == global.PlatformTenantID {
		db.AddError(ErrTenantRequired)
		return "", false
	}
	return tid, true
}

func tableName(stmt *gorm.Statement) string {
	if stmt.Schema != nil {
		return stmt.Schema.Table
	}
	if stmt.TableExpr != nil {
		fields := strings.Fields(stmt.TableExpr.SQL)
		if len(fields) > 0 {
			return strings.Trim(fields[0], "\"`")
		}
	}
	return stmt.Table
}

// Raw SQL on tenant data needs a trusted, explicit bypass; it is never silently scoped.
func (p *TenantPlugin) raw(db *gorm.DB) {
	if !bypassTenant(db) && touchesTenantTable(db.Statement.SQL.String()) {
		db.AddError(errors.New("raw SQL on tenant tables requires an explicit trusted bypass"))
	}
}

func (p *TenantPlugin) scope(db *gorm.DB) {
	if db.Statement.SQL.Len() > 0 {
		p.raw(db)
		return
	}
	var joined []string
	for _, join := range db.Statement.Joins {
		matches := sqlJoin.FindAllStringSubmatch(join.Name, -1)
		for _, match := range matches {
			if global.IsTenantScopedTable(strings.ToLower(match[1])) {
				alias := match[2]
				if alias == "" {
					alias = match[1]
				}
				joined = append(joined, alias)
			}
		}
		if len(matches) == 0 && touchesTenantTable(join.Name) && !bypassTenant(db) {
			db.AddError(errors.New("unsupported tenant join requires an explicit trusted bypass"))
			return
		}
	}
	scoped := global.IsTenantScopedTable(tableName(db.Statement))
	if !scoped && len(joined) == 0 {
		return
	}
	tid, ok := p.resolveTenant(db)
	if !ok {
		return
	}
	if scoped {
		db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: "tenant_id"}, Value: tid}}})
	}
	for _, alias := range joined {
		db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.Eq{Column: clause.Column{Table: alias, Name: "tenant_id"}, Value: tid}}})
	}
}

func checkTenantValue(db *gorm.DB, value reflect.Value, tid string) error {
	if !value.IsValid() {
		return nil
	}
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			if err := checkTenantValue(db, value.Index(i), tid); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			if key.Kind() == reflect.String && (strings.EqualFold(key.String(), "tenant_id") || strings.EqualFold(key.String(), "TenantId")) {
				v := value.MapIndex(key).Interface()
				if str, ok := v.(string); !ok || str != tid {
					return ErrTenantImmutable
				}
			}
		}
	case reflect.Struct:
		field := value.FieldByName("TenantId")
		if field.IsValid() && field.Kind() == reflect.String && field.String() != "" && field.String() != tid {
			return ErrTenantImmutable
		}
	}
	return nil
}

func (p *TenantPlugin) update(db *gorm.DB) {
	p.scope(db)
	if db.Error != nil || !global.IsTenantScopedTable(tableName(db.Statement)) || bypassTenant(db) {
		return
	}
	tid, ok := p.resolveTenant(db)
	if !ok {
		return
	}
	if err := checkTenantValue(db, reflect.ValueOf(db.Statement.Dest), tid); err != nil {
		db.AddError(err)
	}
	value := reflect.ValueOf(db.Statement.Dest)
	for value.IsValid() && (value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}
	if value.IsValid() && value.Kind() == reflect.Struct {
		field := value.FieldByName("TenantId")
		if field.IsValid() && field.Kind() == reflect.String && field.String() == "" {
			for _, selected := range db.Statement.Selects {
				if selected == "*" || strings.EqualFold(selected, "tenant_id") || strings.EqualFold(selected, "TenantId") {
					db.AddError(ErrTenantImmutable)
				}
			}
		}
	}
}

func (p *TenantPlugin) fill(db *gorm.DB) {
	if !global.IsTenantScopedTable(tableName(db.Statement)) {
		return
	}
	tid, ok := p.resolveTenant(db)
	if !ok {
		return
	}
	if err := checkTenantValue(db, reflect.ValueOf(db.Statement.Dest), tid); err != nil {
		db.AddError(err)
		return
	}
	db.Statement.SetColumn("tenant_id", tid)
	// GORM Save may fall back to an upsert. The conflicting row must belong to this tenant.
	if c, ok := db.Statement.Clauses["ON CONFLICT"]; ok {
		if conflict, ok := c.Expression.(clause.OnConflict); ok && !conflict.DoNothing {
			for _, assignment := range conflict.DoUpdates {
				if strings.EqualFold(assignment.Column.Name, "tenant_id") && assignment.Value != tid {
					db.AddError(ErrTenantImmutable)
					return
				}
			}
			conflict.Where.Exprs = append(conflict.Where.Exprs, clause.Eq{Column: clause.Column{Table: db.Statement.Table, Name: "tenant_id"}, Value: tid})
			db.Statement.AddClause(conflict)
		}
	}
}
