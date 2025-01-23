package auth

import (
	"context"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

type AuthService struct {
	// 客户端信息
	models.Client
	// 配置信息
	manage.Config
}

func NewAuthService(auth *AuthService) *server.Server {
	// 1. 创建 OAuth 管理器
	manager := manage.NewDefaultManager()
	// 1.1 设置令牌存储（使用内存存储）
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	// 1.2 设置访问令牌生成器
	manager.MapAccessGenerate(generates.NewAccessGenerate())
	// 1.3 设置客户端存储
	clientStore := store.NewClientStore()
	clientStore.Set(auth.Client.ID, &auth.Client)
	// clientStore.Set("default_client", &models.Client{
	// 	ID:     "default_client", // 默认客户端 ID
	// 	Secret: "default_secret", // 默认客户端密钥
	// 	Domain: "http://localhost:6020",
	// })
	manager.MapClientStorage(clientStore)
	// 1.4 设置 Token 过期时间
	// manager.SetAuthorizeCodeTokenCfg(&manage.Config{
	// 	AccessTokenExp:    time.Hour * 1,      // Access Token 过期时间：1 小时
	// 	RefreshTokenExp:   time.Hour * 24 * 7, // Refresh Token 过期时间：7 天
	// 	IsGenerateRefresh: true,               // 是否生成 Refresh Token
	// })
	manager.SetAuthorizeCodeTokenCfg(&auth.Config)
	manager.SetPasswordTokenCfg(&auth.Config) //  设置密码模式的 Token 过期时间
	// manager.MustTokenStorage(store.NewMemoryTokenStore())// 基于内存创建令牌存储实例

	oauthServer := server.NewServer(server.NewConfig(), manager)
	oauthServer.SetClientInfoHandler(server.ClientFormHandler)                          // 从表单中获取客户端信息
	oauthServer.SetAllowGetAccessRequest(true)                                          // 允许 GET 请求获取令牌
	oauthServer.SetAllowedResponseType(oauth2.ResponseType(oauth2.PasswordCredentials)) // 当前仅允许密码登录

	// 2.1 设置密码授权处理器
	oauthServer.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		// TIP：这里默认已完成用户验证，直接返回用户ID
		return username, nil
	})

	return oauthServer
}
