package initialize

import (
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/pkg/utils/cache"
	"go.uber.org/zap"
)

// InitCache 初始化缓存
func InitCache() {
	cacheType := cache.CacheTypeLocal
	if global.CONFIG.System.UseRedis {
		cacheType = cache.CacheTypeRedis
	}
	addr := fmt.Sprintf("%s:%d", global.CONFIG.Redis.Host, global.CONFIG.Redis.Port)
	cache, err := cache.NewCache(cacheType, addr, global.CONFIG.Redis.Password, global.CONFIG.Redis.DB)
	if err != nil {
		global.LOG.Error("初始化缓存失败", zap.Error(err))
		panic(err)
	}
	global.Cache = cache
}
