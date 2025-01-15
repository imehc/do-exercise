package initialize

import (
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
	cache, err := cache.NewCache(cacheType, global.CONFIG.Redis.Addr, global.CONFIG.Redis.Password, global.CONFIG.Redis.DB)
	if err != nil {
		global.LOG.Error("初始化缓存失败", zap.Error(err))
		panic(err)
	}
	global.Cache = cache
}
