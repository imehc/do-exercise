package util

import (
	"sync"

	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
)

// apiRouteRow 仅取 sys_api 的 method/path 两列。
// 不直接引用 model/system 的 SysApi：model/system 又依赖 util（sys-user.go），
// 会形成 model/system → util → model/system 的导入环。本地声明轻量 struct 即可。
type apiRouteRow struct {
	Method string
	Path   string
}

// 禁用 API 的进程内缓存：以 "METHOD path" 为键。
// sys_api.disabled 是管理界面可切换的接口开关，启用后该接口对所有人（含超管）403。
// 单独一张表 + 内存缓存，避免每个受保护请求都回查 DB。
var (
	disabledApiMu   sync.RWMutex
	disabledApiSet  = map[string]struct{}{}
	disabledApiOnce sync.Once
)

// LoadDisabledApis 从 sys_api 表加载被禁用接口到内存缓存。
// 首次由中间件懒加载；SysApi.Update 变更 disabled 后调用 ReloadDisabledApis 刷新。
func LoadDisabledApis(db *gorm.DB) {
	if db == nil {
		return
	}
	var apis []apiRouteRow
	if err := db.Model(&apiRouteRow{}).Table("sys_api").Select("method", "path").Where("disabled = ?", true).Find(&apis).Error; err != nil {
		return
	}
	disabledApiMu.Lock()
	defer disabledApiMu.Unlock()
	disabledApiSet = make(map[string]struct{}, len(apis))
	for _, api := range apis {
		disabledApiSet[api.Method+" "+api.Path] = struct{}{}
	}
}

// ReloadDisabledApis 重载禁用接口缓存（接口被禁用/启用后调用）。
// 仅 Refresh 语义，不重置懒加载状态；与 LoadDisabledApis 共用实现。
func ReloadDisabledApis() {
	LoadDisabledApis(global.DB)
}

// IsApiDisabled 判断指定接口是否被禁用。
// 首次调用懒加载缓存（依赖 global.DB 已初始化）。
func IsApiDisabled(method, path string) bool {
	disabledApiOnce.Do(func() {
		LoadDisabledApis(global.DB)
	})
	disabledApiMu.RLock()
	defer disabledApiMu.RUnlock()
	_, ok := disabledApiSet[method+" "+path]
	return ok
}
