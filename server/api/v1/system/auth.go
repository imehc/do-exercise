package system

import (
	"context"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/utils"
	"go.uber.org/zap"
)

const (
	PasswordGrantType     = "password"
	RefreshTokenGrantType = "refresh_token"
)

type AuthApi struct{}

// @summary 登录
// @description 使用用户名和密码登录
// @tags auth
// @accept json
// @produce json
// @param request body LoginRequest true "登录请求"
// @success 200 {object} TokenResponse "登录成功"
// @failure 400 {string} string "登录失败"
// @router /login [post]
// @id login
func (a *AuthApi) Login(ctx *gin.Context) {
	req := request.Login{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}

	if ok := global.CAPTCHA_STORE.Verify(req.CaptchaId, req.Captcha, true); ok {
		enable, err := authService.Login(req)
		if err != nil {
			global.LOG.Error("登录失败", zap.Error(err))
			response.BadRequest(err.Error(), ctx)
			return
		}
		if !enable {
			global.LOG.Error("禁止登录", zap.String("username", req.Username))
			response.BadRequest("此用户已禁止登录", ctx)
			return
		}

		// 设置默认值
		// TODO: 处理客户端信息
		values := url.Values{}
		values.Set("grant_type", PasswordGrantType)                  // 默认授权类型
		values.Set("client_id", global.CONFIG.Auth.ClientID)         // 默认客户端 ID
		values.Set("client_secret", global.CONFIG.Auth.ClientSecret) // 默认客户端密钥
		values.Set("username", req.Username)                         // 用户名
		values.Set("password", req.Password)
		// 将值写入请求表单
		ctx.Request.Form = values // 密码
		// 处理令牌请求
		err = global.OAUTH_SERVER.HandleTokenRequest(ctx.Writer, ctx.Request)
		if err != nil {
			response.BadRequest(err.Error(), ctx)
			return
		}
		return
	}
	response.BadRequest("验证码错误", ctx)
}

// @summary 登出
// @description 退出登录
// @tags auth
// @accept json
// @produce json
// @security bearerauth
// @success 204 {object} nil "登出成功"
// @failure 400 {string} string "登出失败"
// @failure 401 {string} string "未授权"
// @router /logout [get]
// @id logout
func (a *AuthApi) Logout(ctx *gin.Context) {
	// claims := c.MustGet("claims").(system.Claims)
	// TODO: 从redis中删除对应信息
	accessToken := strings.Split(ctx.Request.Header.Get("Authorization"), " ")[1]
	global.OAUTH_SERVER.Manager.RemoveRefreshToken(context.Background(), accessToken)
	response.NoContent(ctx)
}

// @summary 刷新令牌
// @description 使用refresh_token刷新access_token
// @tags auth
// @accept json
// @produce json
// @security bearerauth
// @param refresh_token query string true "刷新令牌"
// @success 200 {object} TokenResponse "刷新令牌成功"
// @failure 400 {string} string "刷新令牌失败"
// @router /refresh_token [get]
// @id refresh_token
func (a *AuthApi) RefreshToken(ctx *gin.Context) {
	req := request.RefreshToken{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(utils.GetValidMsg(req, err), ctx)
		return
	}
	// 设置默认值
	values := url.Values{}
	values.Set("grant_type", RefreshTokenGrantType)              // 授权类型为刷新 Token
	values.Set("client_id", global.CONFIG.Auth.ClientID)         // 默认客户端 ID
	values.Set("client_secret", global.CONFIG.Auth.ClientSecret) // 默认客户端密钥
	values.Set("refresh_token", req.RefreshToken)
	ctx.Request.Form = values

	// 自定义错误返回响应
	global.OAUTH_SERVER.SetResponseErrorHandler(func(re *errors.Response) {
		re.Description = ""
		re.Error = nil
		ctx.String(response.BAD_REQUEST, "refresh_token is invalid")
	})

	// 处理令牌请求
	err := global.OAUTH_SERVER.HandleTokenRequest(ctx.Writer, ctx.Request)
	if err != nil {
		response.BadRequest(err.Error(), ctx)
		return
	}
}
