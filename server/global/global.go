package global

import (
	"github.com/imehc/do-exercise/server/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	// Config 全局配置
	Config *config.Config
	// DB 全局数据库连接
	DB *gorm.DB
	// Redis 全局Redis客户端
	Redis *redis.Client
)
