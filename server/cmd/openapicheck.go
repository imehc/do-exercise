package cmd

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/imehc/do-exercise/server/internal"
	"github.com/imehc/do-exercise/server/router"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// openapiCheckCmd 校验 openapi.yaml 与 gin 实际注册路由的一致性（漂移检查）。
// 两侧都以 "METHOD /path" 归一化后比对：任何差异都会以非零退出码结束，
// 供 CI / Makefile 使用，防止接口文档与代码脱节。
// infraRouteAllowlist 内部基础设施端点，不属于业务 API 文档范围，
// 允许代码注册但无需出现在 openapi.yaml（健康检查、SSE 事件流等）。
var infraRouteAllowlist = map[string]bool{
	"GET /health": true,
	"GET /sse":    true,
}

var openapiCheckCmd = &cobra.Command{
	Use:   "openapicheck",
	Short: "检查 openapi.yaml 与运行时路由漂移",
	Long: `读取 openapi.yaml，用 gin 实际注册的路由做双向比对：
- 代码已注册、spec 缺失 → 文档欠录（错误）
- spec 存在、代码未注册 → 文档悬空（错误）`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.InitConfig(configFile)
		r := router.New()
		if r == nil {
			fmt.Println("构建路由失败")
			os.Exit(1)
		}

		routes := make(map[string]bool)
		for _, route := range r.Routes() {
			key := route.Method + " " + ginParamToSpec(route.Path)
			if infraRouteAllowlist[key] {
				continue
			}
			routes[key] = true
		}

		spec, err := loadSpecPaths("openapi.yaml")
		if err != nil {
			fmt.Printf("读取 openapi.yaml 失败: %v\n", err)
			os.Exit(1)
		}

		missing := difference(routes, spec)
		extra := difference(spec, routes)

		if len(missing) > 0 || len(extra) > 0 {
			if len(missing) > 0 {
				sort.Strings(missing)
				fmt.Printf("✗ [spec 缺失] 已注册但未写入 openapi.yaml（%d）:\n", len(missing))
				for _, k := range missing {
					fmt.Printf("    %s\n", k)
				}
			}
			if len(extra) > 0 {
				sort.Strings(extra)
				fmt.Printf("✗ [spec 悬空] openapi.yaml 有但未注册路由（%d）:\n", len(extra))
				for _, k := range extra {
					fmt.Printf("    %s\n", k)
				}
			}
			os.Exit(1)
		}

		fmt.Printf("✓ openapi.yaml 与路由一致（共 %d 条）\n", len(routes))
	},
}

// ginParamToSpec 把 gin 的 :param 路径段转成 openapi 的 {param} 形式，
// 保证两侧可比较（/auth/:id 与 /auth/{id} 视为同一路由）。
func ginParamToSpec(path string) string {
	return regexp.MustCompile(`:([A-Za-z0-9_]+)`).ReplaceAllString(path, "{$1}")
}

// loadSpecPaths 解析 openapi.yaml 顶层 paths，返回 "METHOD /path" 集合。
func loadSpecPaths(file string) (map[string]bool, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var doc struct {
		Paths map[string]map[string]any `yaml:"paths"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	result := make(map[string]bool)
	methods := []string{"get", "post", "put", "delete", "patch", "head", "options"}
	for path, ops := range doc.Paths {
		for _, method := range methods {
			if _, ok := ops[method]; ok {
				result[methodToUpper(method)+" "+path] = true
			}
		}
	}
	return result, nil
}

func methodToUpper(m string) string {
	upper := map[string]string{
		"get": "GET", "post": "POST", "put": "PUT", "delete": "DELETE",
		"patch": "PATCH", "head": "HEAD", "options": "OPTIONS",
	}
	return upper[m]
}

// difference 返回 a 中有而 b 中没有的键。
func difference(a, b map[string]bool) []string {
	var out []string
	for k := range a {
		if !b[k] {
			out = append(out, k)
		}
	}
	return out
}

func init() {
	rootCmd.AddCommand(openapiCheckCmd)
}
