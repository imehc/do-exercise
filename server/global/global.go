package global

import (
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/imehc/do-exercise/server/config"
	"github.com/imehc/do-exercise/server/pkg/utils/cache"
	"github.com/mojocn/base64Captcha"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var (
	VP                  *viper.Viper            // 配置
	CONFIG              *config.Config          // 配置
	OAUTH_SERVER        *server.Server          // 认证服务器
	LOG                 *zap.Logger             // 日志
	Cache               cache.Cache             // 缓存
	Concurrency_Control = &singleflight.Group{} // 并发控制
	CAPTCHA_STORE       base64Captcha.Store     // 验证码存储
)
