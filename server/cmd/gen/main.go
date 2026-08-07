package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// multiFlag 收集重复的 -field 参数
type multiFlag struct {
	fields []Field
}

func (m *multiFlag) String() string {
	names := make([]string, 0, len(m.fields))
	for _, f := range m.fields {
		names = append(names, f.Name)
	}
	return strings.Join(names, ",")
}

func (m *multiFlag) Set(v string) error {
	f, err := parseField(v)
	if err != nil {
		return err
	}
	m.fields = append(m.fields, f)
	return nil
}

func main() {
	name := flag.String("name", "", "实体结构体名（如 SysDevice）")
	desc := flag.String("desc", "", "实体描述（如 设备）")
	path := flag.String("path", "", "路由分组路径（默认取实体名短横线形式）")
	dir := flag.String("dir", ".", "server 根目录")
	var fields multiFlag
	flag.Var(&fields, "field", "字段，格式 name:type[:注释]，可多次指定；type: string/text/bool/int/uint/int64/float/time/id_array")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: go run ./cmd/gen -name=SysDevice -desc=设备 -path=devices -field=name:string:设备名称 -field=enabled:bool:是否启用\n")
		flag.PrintDefaults()
	}

	if *name == "" {
		flag.Usage()
		os.Exit(1)
	}

	spec := &Spec{
		Name:   *name,
		Desc:   *desc,
		Path:   *path,
		Fields: fields.fields,
	}
	if spec.Path == "" {
		spec.Path = spec.kebab()
	}

	// 生成六件套
	if err := spec.writeFiles(*dir); err != nil {
		fmt.Fprintln(os.Stderr, "生成失败:", err)
		os.Exit(1)
	}
	// 注册到 enter 文件、路由、迁移
	if err := spec.register(*dir); err != nil {
		fmt.Fprintln(os.Stderr, "注册失败:", err)
		os.Exit(1)
	}

	fmt.Printf("\n生成完成。请在 server/ 下执行 go build ./... 验证，并检查是否需要对生成文件补充逻辑。\n")
}
