package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type patchItem struct {
	needle string
	add    string
	exists string // 幂等判据：已包含该字符串（任意位置）则跳过；空则用 add 本身
}

// insertBefore 在 needle 之前插入 add（含换行）
func insertBefore(content, needle, add string) (string, bool) {
	idx := strings.Index(content, needle)
	if idx < 0 {
		return content, false
	}
	return content[:idx] + add + content[idx:], true
}

// patchFile 对文件按顺序执行多处幂等插入
func patchFile(path string, items []patchItem) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取 %s 失败: %w", path, err)
	}
	text := string(content)
	changed := false
	for _, it := range items {
		existsCheck := it.exists
		if existsCheck == "" {
			existsCheck = it.add
		}
		if strings.Contains(text, existsCheck) {
			continue
		}
		var ok bool
		text, ok = insertBefore(text, it.needle, it.add)
		if ok {
			changed = true
		}
	}
	if changed {
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			return err
		}
		fmt.Printf("已更新 %s\n", path)
	}
	return nil
}

// register 将实体注册到各 enter 文件、路由与迁移列表
func (spec *Spec) register(dir string) error {
	svc := spec.Name + "Service"
	api := spec.Name + "Api"
	rt := spec.Name + "Router"

	// 1. service/system/enter.go：ServiceGroup 增加字段
	p := filepath.Join(dir, "service", "system", "enter.go")
	if err := patchFile(p, []patchItem{
		{needle: "}", add: "\t" + svc + "\n"}, // 追加到 ServiceGroup 最后一个字段后
	}); err != nil {
		return err
	}

	// 2. api/v1/system/enter.go：ApiGroup 增加字段 + service 变量
	p = filepath.Join(dir, "api", "v1", "system", "enter.go")
	svcAssign := "\t" + spec.serviceVar() + " = service.ServiceGroupApp.SystemServiceGroup." + svc + " // " + spec.Desc + "服务\n"
	if err := patchFile(p, []patchItem{
		{needle: "}", add: "\t" + api + "\n", exists: "\t" + api + "\n"},
		{needle: ")", add: svcAssign, exists: spec.serviceVar()},
	}); err != nil {
		return err
	}

	// 3. router/system/enter.go：RouterGroup 增加字段 + api 变量
	p = filepath.Join(dir, "router", "system", "enter.go")
	apiAssign := "\t" + spec.apiVar() + " = api.ApiGroupApp.SystemApiGroup." + api + " // " + spec.Desc + "接口\n"
	if err := patchFile(p, []patchItem{
		{needle: "}", add: "\t" + rt + "\n", exists: "\t" + rt + "\n"},
		{needle: ")", add: apiAssign, exists: spec.apiVar()},
	}); err != nil {
		return err
	}

	// 4. router/router.go：挂载路由
	p = filepath.Join(dir, "router", "router.go")
	if err := patchFile(p, []patchItem{
		{needle: "system.InitSysJobRouter(protected)", add: "\t\t" + "system.Init" + rt + "(protected)\n"},
	}); err != nil {
		return err
	}

	// 5. internal/gorm.go：AutoMigrate 注册模型
	p = filepath.Join(dir, "internal", "gorm.go")
	if err := patchFile(p, []patchItem{
		{needle: "gormadapter.CasbinRule{},", add: "\t\t\tsystem." + spec.Name + "{},\n"},
	}); err != nil {
		return err
	}

	return nil
}
