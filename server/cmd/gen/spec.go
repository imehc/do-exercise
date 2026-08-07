package main

import (
	"fmt"
	"strings"
)

// Field 描述生成实体的一列字段
type Field struct {
	Name    string // 蛇形字段名，如 "device_name"
	Pascal  string // 结构体字段名，如 "DeviceName"
	GoType  string // Go 类型
	Kind    string // string/bool/int/time/text/...
	Comment string // 注释（可选）
}

// Spec 由命令行参数解析出的生成规格
type Spec struct {
	Name   string  // 实体结构体名，如 "SysDevice"
	Desc   string  // 实体描述
	Path   string  // 路由分组路径
	Fields []Field // 自定义字段
}

func parseField(raw string) (Field, error) {
	parts := strings.Split(raw, ":")
	if len(parts) < 2 {
		return Field{}, fmt.Errorf("字段格式应为 name:type[:注释]，收到：%q", raw)
	}
	name := parts[0]
	kind := parts[1]
	comment := ""
	if len(parts) >= 3 {
		comment = parts[2]
	}

	field := Field{
		Name:    name,
		Pascal:  pascal(name),
		Kind:    kind,
		Comment: comment,
	}

	switch kind {
	case "string":
		field.GoType = "string"
	case "text":
		field.GoType = "string"
	case "bool":
		field.GoType = "bool"
	case "int":
		field.GoType = "int"
	case "uint":
		field.GoType = "uint"
	case "int64":
		field.GoType = "int64"
	case "float":
		field.GoType = "float64"
	case "time":
		field.GoType = "time.Time"
	case "id_array":
		field.GoType = "[]uint"
	default:
		return Field{}, fmt.Errorf("不支持的类型 %q（可选：string/text/bool/int/uint/int64/float/time/id_array）", kind)
	}
	return field, nil
}

// KeyField 是否为时间/删除等由骨架自动生成的字段
func (spec *Spec) structName() string    { return spec.Name }
func (spec *Spec) serviceStruct() string { return spec.Name + "Service" }
func (spec *Spec) apiStruct() string     { return spec.Name + "Api" }
func (spec *Spec) routerStruct() string  { return spec.Name + "Router" }
func (spec *Spec) createReq() string     { return "Create" + spec.Name + "Req" }
func (spec *Spec) updateReq() string     { return "Update" + spec.Name + "Req" }
func (spec *Spec) resp() string          { return spec.Name + "Resp" }
func (spec *Spec) kebab() string         { return toKebab(spec.Name) }
func (spec *Spec) fileName() string      { return toKebab(spec.Name) + ".go" }
func (spec *Spec) serviceVar() string    { return lowerFirst(spec.Name) + "Service" }
func (spec *Spec) apiVar() string        { return lowerFirst(spec.Name) + "Api" }

func pascal(name string) string {
	parts := strings.Split(strings.ToLower(name), "_")
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}

func toKebab(name string) string {
	var b strings.Builder
	runes := []rune(name)
	for i, r := range runes {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('-')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}

func lowerFirst(name string) string {
	if name == "" {
		return name
	}
	return strings.ToLower(name[:1]) + name[1:]
}
