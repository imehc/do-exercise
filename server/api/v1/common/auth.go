package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type AuthApi struct{}

// buildToken 生成带租户信息的访问令牌
func (s *AuthApi) buildToken(ctx *gin.Context, user *system.SysUser) (*common.Token, error) {
	baseConf := util.Token{
		UserId:   user.UserId,
		Username: user.Username,
		RoleIds: lo.Map(user.Roles, func(item system.SysRole, index int) uint {
			return item.Id
		}),
		TenantId:           user.TenantId,
		ExpireTime:         util.AuthAccessExpire(),
		RefreshExpireTime:  util.AuthRefreshExpire(),
		Disabled:           false, // TODO: 在数据库中添加字段获取禁用状态
		CreatedTime:        time.Now(),
		MustChangePassword: user.MustChangePassword,
		IsSuperAdmin:       user.IsSuperAdmin,
	}
	token, err := baseConf.GenerateToken()
	if err != nil {
		return nil, err
	}
	return token, nil
}

// respondLoginResult 统一返回登录/选租户/切换租户响应（token + 可用租户列表）
func (s *AuthApi) respondLoginResult(ctx *gin.Context, token *common.Token, available []response.TenantOption) {
	response.Success(ctx, response.LoginResult{
		Token:            token,
		AvailableTenants: available,
	})
}

// Login 登录
func (s *AuthApi) Login(ctx *gin.Context) {
	var req common.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if ok := global.Captcha.Verify(req.CaptchaId, req.Captcha, true); !ok {
		response.BadRequest(ctx, "invalidCaptcha")
		return
	}

	password, err := shared.RSACrypto.DecryptWithKey(req.PublicKey, req.Password, true)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// 登录爆破防护：按IP与用户名计数，失败次数越多等待越久，达到阈值后硬锁
	ip := ctx.ClientIP()
	locked, delay := loginPenalty(ip, req.Username)
	if locked {
		response.BadRequest(ctx, "loginLocked")
		return
	}
	if delay > 0 {
		time.Sleep(delay)
	}

	user, options, err := authService.Login(util.DB(ctx), common.Login{
		Username: req.Username,
		Password: password,
	})
	if err != nil {
		if err.Error() == "requiresTenantSelection" {
			// 认证已通过，仅需选择租户：签发一次性登录会话供 select_tenant 使用
			clearLoginFailures(ip, req.Username)
			sessionId, serr := s.createLoginSession(ctx, req.Username, options)
			if serr != nil {
				response.ServerError(ctx)
				return
			}
			response.Success(ctx, response.LoginResult{
				RequiresTenantSelection: true,
				LoginSessionId:          sessionId,
				AvailableTenants:        options,
			})
			return
		}
		registerLoginFailure(ip, req.Username)
		response.BadRequest(ctx, err.Error())
		return
	}

	clearLoginFailures(ip, req.Username)
	token, terr := s.buildToken(ctx, user)
	if terr != nil {
		response.ServerError(ctx)
		return
	}
	s.respondLoginResult(ctx, token, options)
}

// RefreshToken 刷新token。
// refresh token 走请求体而不是查询串，避免进入访问日志、浏览器历史与 Referer；
// 每次刷新都会轮转会刷新令牌本身（见 util/token.go 的 RefreshToken）。
func (s *AuthApi) RefreshToken(ctx *gin.Context) {
	type RefreshTokenReq struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	req := RefreshTokenReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "emptyRefreshToken")
		return
	}
	baseConf := util.Token{
		ExpireTime:        util.AuthAccessExpire(),
		RefreshExpireTime: util.AuthRefreshExpire(),
	}
	token, err := baseConf.RefreshToken(req.RefreshToken)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, token)
}

// PublicKey 获取公钥
func (s *AuthApi) PublicKey(ctx *gin.Context) {
	publicKey, error := shared.RSACrypto.GetRandomKeyPair()
	if error != nil {
		response.ServerError(ctx)
		return
	}
	response.Success(ctx, publicKey)
}

// GetCaptcha 获取验证码
func (s *AuthApi) GetCaptcha(ctx *gin.Context) {
	id, b64s, _, err := global.Captcha.Generate()
	if err != nil {
		global.Log.Error("验证码获取失败!", zap.Error(err))
		response.BadRequest(ctx, "captchaError")
		return
	}
	response.Success(ctx, common.Captcha{
		CaptchaId: id,
		PicPath:   b64s,
	})
}

// AvailableTenants 可用租户列表（兼容保留）。
// 前端登录页选择器已移除，该接口仅在多租户模式下返回启用的业务租户（不含平台保留租户）。
func (s *AuthApi) AvailableTenants(ctx *gin.Context) {
	rep := response.AvailableTenants{
		Mode:    global.Config.Tenant.Mode,
		Tenants: []response.TenantOption{},
	}
	if global.Config.Tenant.IsMulti() {
		tenants, err := tenantPublicService.ListEnabled(util.DB(ctx))
		if err != nil {
			response.ServerError(ctx)
			return
		}
		rep.Tenants = tenants
	}
	response.Success(ctx, rep)
}

// ResetPassword 忘记密码
func (s *AuthApi) ResetPassword(ctx *gin.Context) {
	iRedis := global.Redis
	context := context.Background()

	req := &request.UserResetPasswordReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	user, err := userService.FindUserByEmail(util.DB(ctx), req.Email)
	if err != nil {
		// 未绑定的邮箱与“验证码错误”返回同一个响应，避免枚举账号存在性。
		response.BadRequest(ctx, "captchaError")
		return
	}

	password, err := shared.RSACrypto.DecryptWithKey(req.PublicKey, req.Password, true)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	var userApi UserApi
	cache, err := userApi.getEmailCache(context, iRedis, ForgotPasswordPrefix, req.Email)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	defer func() {
		_ = userApi.clearEmailCache(context, iRedis, ForgotPasswordPrefix, req.Email)
	}()

	if cache.Code != req.Code || cache.UserId != user.UserId {
		response.BadRequest(ctx, "captchaError")
		return
	}

	if err := authService.ResetPassword(util.DB(ctx), user.UserId, request.UserResetPasswordReq{
		Password: password,
	}); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// SendResetPasswordCode 发送忘记密码邮箱验证码
func (s *AuthApi) SendResetPasswordCode(ctx *gin.Context) {
	req := &request.Email{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	user, err := userService.FindUserByEmail(util.DB(ctx), req.Email)
	if err != nil {
		// 未绑定的邮箱与“验证码已发送”返回同一个响应，避免枚举账号存在性。
		// 不真正发信，但保持一致的 204，客户端无法区分。
		response.NoContent(ctx)
		return
	}

	var userApi UserApi
	userApi.sendEmailCode(ctx, ForgotPasswordPrefix, &util.EmailData{
		EmailTitle:       "验证码",
		VerificationType: "重置密码",
		GreetingText:     "感谢您使用我们的服务，请使用以下验证码完成密码重置：",
	}, user.UserId)
}

// Logout 退出登录
func (s *AuthApi) Logout(ctx *gin.Context) {
	accessToken := ctx.GetHeader("Authorization")
	userId := ctx.MustGet("userId").(string)

	if err := authService.Logout(userId, accessToken); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.NoContent(ctx)
}

const (
	loginSessionPrefix = "loginSession:"
	loginSessionTTL    = 5 * time.Minute
)

// createLoginSession 为“多租户待选择”的登录签发一次性会话
func (s *AuthApi) createLoginSession(ctx *gin.Context, username string, options []response.TenantOption) (string, error) {
	id, err := util.Uuid()
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(common.LoginSession{
		Username: username,
		Tenants: lo.Map(options, func(o response.TenantOption, _ int) string {
			return o.TenantId
		}),
	})
	if err != nil {
		return "", err
	}
	if err := global.Redis.Set(context.Background(), fmt.Sprintf("%s%s", loginSessionPrefix, id), jsonBytes, loginSessionTTL).Err(); err != nil {
		return "", err
	}
	return id, nil
}

// getLoginSession 读取登录会话；不存在或已失效返回 loginSessionExpired
func (s *AuthApi) getLoginSession(sessionId string) (*common.LoginSession, error) {
	data, err := global.Redis.Get(context.Background(), fmt.Sprintf("%s%s", loginSessionPrefix, sessionId)).Result()
	if err != nil {
		return nil, errors.New("loginSessionExpired")
	}
	var session common.LoginSession
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, errors.New("loginSessionExpired")
	}
	return &session, nil
}

// destroyLoginSession 销毁登录会话（一次性消费）
func (s *AuthApi) destroyLoginSession(sessionId string) error {
	return global.Redis.Del(context.Background(), fmt.Sprintf("%s%s", loginSessionPrefix, sessionId)).Err()
}

// SelectTenant 多租户登录选择租户。登录会话一次性使用且只能进入会话认证过的租户，
// 防止会话复用或横向切换到未获授权的租户。
func (s *AuthApi) SelectTenant(ctx *gin.Context) {
	var req common.SelectTenantReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	session, err := s.getLoginSession(req.LoginSessionId)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	// 无论成败都销毁会话，杜绝重用
	_ = s.destroyLoginSession(req.LoginSessionId)

	if !lo.Contains(session.Tenants, req.TenantId) {
		response.BadRequest(ctx, "invalidTenant")
		return
	}

	user, err := authService.EnterTenant(util.DB(ctx), req.TenantId, session.Username)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	token, terr := s.buildToken(ctx, user)
	if terr != nil {
		response.ServerError(ctx)
		return
	}
	options, _ := tenantPublicService.ListEnabledForUsername(util.DB(ctx), session.Username)
	s.respondLoginResult(ctx, token, options)
}

// SwitchTenant 登录后切换租户。重新签发目标租户的 token，并吊销当前会话的旧凭据。
func (s *AuthApi) SwitchTenant(ctx *gin.Context) {
	var req common.SwitchTenantReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	username := ctx.GetString("username")
	user, err := authService.EnterTenant(util.BypassTenantDB(ctx), req.TenantId, username)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	token, terr := s.buildToken(ctx, user)
	if terr != nil {
		response.ServerError(ctx)
		return
	}
	// 吊销当前租户的会话凭据，避免旧租户凭据继续残留
	_ = authService.Logout(ctx.GetString("userId"), ctx.GetString("accessToken"))

	options, _ := tenantPublicService.ListEnabledForUsername(util.BypassTenantDB(ctx), username)
	s.respondLoginResult(ctx, token, options)
}

// MyTenants 当前用户的可用租户列表（租户切换器用）。平台保留租户不在此列。
func (s *AuthApi) MyTenants(ctx *gin.Context) {
	username := ctx.GetString("username")
	options, err := tenantPublicService.ListEnabledForUsername(util.BypassTenantDB(ctx), username)
	if err != nil {
		response.ServerError(ctx)
		return
	}
	response.Success(ctx, options)
}

// LoginWithEmail 使用邮件登录
func (s *AuthApi) LoginWithEmail(ctx *gin.Context) {
	context := context.Background()
	iRedis := global.Redis

	var req common.LoginWithEmailReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	var userApi UserApi
	cache, err := userApi.getEmailCache(context, iRedis, LoginWithEmailPrefix, req.Email)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	defer func() {
		_ = userApi.clearEmailCache(context, iRedis, LoginWithEmailPrefix, req.Email)
	}()

	user, err := userService.FindUserByEmail(util.DB(ctx), req.Email)
	if err != nil {
		// 未绑定的邮箱与“验证码错误”返回同一个响应，避免枚举账号存在性。
		response.BadRequest(ctx, "captchaError")
		return
	}

	if cache.Code != req.Code || cache.UserId != user.UserId {
		response.BadRequest(ctx, "captchaError")
		return
	}

	token, err := s.buildToken(ctx, user)
	if err != nil {
		response.ServerError(ctx)
		return
	}
	response.Success(ctx, token)
}

// SendLoginWithEmailCode 发送使用邮箱登录验证码
func (s *AuthApi) SendLoginWithEmailCode(ctx *gin.Context) {
	req := &request.Email{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	user, err := userService.FindUserByEmail(util.DB(ctx), req.Email)
	if err != nil {
		// 未绑定的邮箱与“验证码已发送”返回同一个响应，避免枚举账号存在性。
		// 不真正发信，但保持一致的 204，客户端无法区分。
		response.NoContent(ctx)
		return
	}

	var userApi UserApi
	userApi.sendEmailCode(ctx, LoginWithEmailPrefix, &util.EmailData{
		EmailTitle:       "验证码",
		VerificationType: "邮箱登录",
		GreetingText:     "感谢您使用我们的服务，请使用以下验证码完成登录：",
	}, user.UserId)
}
