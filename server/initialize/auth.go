package initialize

import (
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/pkg/auth"
	"github.com/imehc/do-exercise/server/pkg/utils"
	"go.uber.org/zap"
)

// InitAuth 初始化认证服务
func InitAuth() {
	accessExpire, err := utils.ParseDurationString(global.CONFIG.Auth.AccessExpire)
	if err != nil {
		global.LOG.Error("auth config error", zap.Error(err))
		panic(err)
	}
	refreshExpire, err2 := utils.ParseDurationString(global.CONFIG.Auth.RefreshExpire)
	if err2 != nil {
		global.LOG.Error("auth config error", zap.Error(err2))
		panic(err2)
	}

	config := NewAuthConfig(
		WithClient(models.Client{
			ID:     global.CONFIG.Auth.ClientID,     // 默认客户端 ID
			Secret: global.CONFIG.Auth.ClientSecret, // 默认客户端密钥
			Domain: global.CONFIG.Auth.Domain,
		}),
		WithConfig(manage.Config{
			AccessTokenExp:    accessExpire,  // Access Token 过期时间：2 小时
			RefreshTokenExp:   refreshExpire, // Refresh Token 过期时间：7 天
			IsGenerateRefresh: true,          // 是否生成 Refresh Token
		}),
	)

	global.OAUTH_SERVER = auth.NewAuthService(config)
}

// Option 配置选项
type Option func(*auth.AuthService)

// NewAuthConfig 创建 AuthService 配置
func NewAuthConfig(opts ...Option) *auth.AuthService {
	config := &auth.AuthService{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

// WithClient 设置客户端
func WithClient(client models.Client) Option {
	return func(config *auth.AuthService) {
		config.Client = client
	}
}

// WithConfig 设置配置
func WithConfig(conf manage.Config) Option {
	return func(config *auth.AuthService) {
		config.Config = conf
	}
}
