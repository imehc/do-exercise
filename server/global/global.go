package global

import (
	"github.com/casbin/casbin/v3"
	"github.com/imehc/do-exercise/server/config"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/minio/minio-go/v7"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config   *config.Config         // Config 全局配置
	DB       *gorm.DB               // DB 全局数据库连接
	Redis    *redis.Client          // Redis 全局Redis客户端
	I18      *I18n                  // I18 国际化翻译器
	Log      *zap.Logger            // Log 全局日志
	Captcha  *base64Captcha.Captcha // Captcha 验证码
	Enforcer *casbin.Enforcer       // Enforcer 权限控制器
	Searcher *xdb.Searcher          // Searcher IP2Region查询器
	Oss      *minio.Client          // Oss OSS对象存储
)
