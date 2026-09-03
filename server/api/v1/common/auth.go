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
func (s *AuthApi) buildToken(ctx *gin.Context, user *system.SysUser, available []response.TenantOption) (*common.Token, error) {
	baseConf := util.Token{
		UserId:   user.UserId,
		Username: user.Username,
		TenantId: user.TenantId,
		AuthorizedTenantIds: lo.Map(available, func(item response.TenantOption, _ int) string {
			return item.TenantId
		}),
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
		Username:   req.Username,
		Password:   password,
		TenantId:   req.TenantId,
		TenantCode: req.TenantCode,
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
	token, terr := s.buildToken(ctx, user, options)
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

// AvailableTenants 返回登录页需要的静态词表。
// 登录页已不再允许匿名枚举租户，因此 tenants 始终为空；认证后的候选租户由登录结果返回。
func (s *AuthApi) AvailableTenants(ctx *gin.Context) {
	rep := response.AvailableTenants{
		Tenants:           []response.TenantOption{},
		PermissionActions: global.MenuPermissionActions,
	}
	response.Success(ctx, rep)
}

// ResetPassword 忘记密码。
//
// 邮箱在多租户下不唯一，验证码只能证明「请求者拥有这个邮箱」，证明不了他要改哪个租户
// 的口令。因此候选账号多于一个时不动任何数据，返回候选租户要求显式指定 tenant_id
// ——一次请求重置同邮箱在所有租户下的口令，等于凭一个邮箱横向接管多个租户的账号。
func (s *AuthApi) ResetPassword(ctx *gin.Context) {
	iRedis := global.Redis
	context := context.Background()

	req := &request.UserResetPasswordReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
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

	// 验证码只在拿到终态（重置成功，或校验失败）时消费。
	// 「还需选择租户」不是终态：此时一个字节都没改，若在这里清掉缓存，
	// 用户带着 tenant_id 重试必然失败，只能重新收信。
	consume := true
	defer func() {
		if consume {
			_ = userApi.clearEmailCache(context, iRedis, ForgotPasswordPrefix, req.Email)
		}
	}()

	users, err := s.emailCandidates(ctx, cache, req.Email, req.Code)
	if err != nil {
		// 未绑定的邮箱与「验证码错误」返回同一个响应，避免枚举账号存在性。
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.TenantId != "" {
		users = lo.Filter(users, func(user system.SysUser, _ int) bool {
			return user.TenantId == req.TenantId
		})
		if len(users) != 1 {
			response.BadRequest(ctx, "invalidTenant")
			return
		}
	} else if len(users) > 1 {
		consume = false
		options, _ := tenantPublicService.ListEnabledByIDs(util.DB(ctx), userService.TenantIdsOf(users))
		response.Success(ctx, response.ResetPasswordResult{
			RequiresTenantSelection: true,
			AvailableTenants:        options,
		})
		return
	}

	if err := authService.ResetPassword(util.DB(ctx), users[0].UserId, request.UserResetPasswordReq{
		Password: password,
	}); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, response.ResetPasswordResult{AvailableTenants: []response.TenantOption{}})
}

// emailCandidates 校验邮箱验证码并收敛出候选账号。
//
// 「邮箱查不到」「验证码不对」「验证码不是发给这批账号的」三种情况返回同一个
// captchaError：任何差异都能被用来枚举某个邮箱是否注册过。
//
// 候选集取「当前仍可用的账号」与「发码时绑定的账号」的交集——发码后才绑定同一邮箱的
// 新账号不在集合内，否则抢先绑定邮箱就能蹭到别人已经收到的验证码。
func (s *AuthApi) emailCandidates(ctx *gin.Context, cache *request.EmailCache, email, code string) ([]system.SysUser, error) {
	users, err := userService.FindUsersByEmail(util.DB(ctx), email)
	if err != nil || cache.Code != code {
		return nil, errors.New("captchaError")
	}
	users = lo.Filter(users, func(user system.SysUser, _ int) bool {
		return cache.Allows(user.UserId)
	})
	if len(users) == 0 {
		return nil, errors.New("captchaError")
	}
	return users, nil
}

// SendResetPasswordCode 发送忘记密码邮箱验证码
func (s *AuthApi) SendResetPasswordCode(ctx *gin.Context) {
	req := &request.Email{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	users, err := userService.FindUsersByEmail(util.DB(ctx), req.Email)
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
	}, userService.UserIdsOf(users))
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

// createLoginSession 为「待选择租户」的口令登录签发一次性会话
func (s *AuthApi) createLoginSession(ctx *gin.Context, username string, options []response.TenantOption) (string, error) {
	return s.saveLoginSession(common.LoginSession{
		Username: username,
		Tenants:  tenantIdsOf(options),
	})
}

// createEmailLoginSession 为「待选择租户」的邮箱登录签发一次性会话。
// 邮箱路径记录候选账号 ID 而非用户名，见 common.LoginSession 的说明。
func (s *AuthApi) createEmailLoginSession(userIds []string, options []response.TenantOption) (string, error) {
	return s.saveLoginSession(common.LoginSession{
		Tenants: tenantIdsOf(options),
		UserIds: userIds,
	})
}

func tenantIdsOf(options []response.TenantOption) []string {
	return lo.Map(options, func(o response.TenantOption, _ int) string { return o.TenantId })
}

func (s *AuthApi) saveLoginSession(session common.LoginSession) (string, error) {
	id, err := util.Uuid()
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(session)
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

// SelectTenant 登录时选择租户。登录会话一次性使用且只能进入会话认证过的租户，
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

	// 邮箱登录的会话按候选账号 ID 收敛，口令登录的会话按用户名收敛。
	var user *system.SysUser
	if len(session.UserIds) > 0 {
		user, err = authService.EnterTenantByUserIds(util.DB(ctx), req.TenantId, session.UserIds)
	} else {
		user, err = authService.EnterTenant(util.DB(ctx), req.TenantId, session.Username)
	}
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	options, _ := tenantPublicService.ListEnabledByIDs(util.DB(ctx), session.Tenants)
	token, terr := s.buildToken(ctx, user, options)
	if terr != nil {
		response.ServerError(ctx)
		return
	}
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
	authorizedTenantIds := ctx.GetStringSlice("authorizedTenantIds")
	if !lo.Contains(authorizedTenantIds, req.TenantId) {
		response.BadRequest(ctx, "invalidTenant")
		return
	}
	password, err := shared.RSACrypto.DecryptWithKey(req.PublicKey, req.Password, true)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	user, err := authService.Authenticate(util.BypassTenantDB(ctx), req.TenantId, username, password)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	options, _ := tenantPublicService.ListEnabledByIDs(util.BypassTenantDB(ctx), authorizedTenantIds)
	token, terr := s.buildToken(ctx, user, options)
	if terr != nil {
		response.ServerError(ctx)
		return
	}
	// 吊销当前租户的会话凭据，避免旧租户凭据继续残留
	_ = authService.Logout(ctx.GetString("userId"), ctx.GetString("accessToken"))

	s.respondLoginResult(ctx, token, options)
}

// MyTenants 当前用户的可用租户列表（租户切换器用）。平台保留租户不在此列。
func (s *AuthApi) MyTenants(ctx *gin.Context) {
	options, err := tenantPublicService.ListEnabledByIDs(
		util.BypassTenantDB(ctx),
		ctx.GetStringSlice("authorizedTenantIds"),
	)
	if err != nil {
		response.ServerError(ctx)
		return
	}
	response.Success(ctx, options)
}

// LoginWithEmail 使用邮箱验证码登录。
//
// 与口令登录同构：邮箱命中多个启用租户时不擅自挑一个，返回 requires_tenant_selection
// 与一次性登录会话，由前端选择后经 select_tenant 完成登录。
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

	// 验证码在这里一次性消费：后续的租户选择由一次性登录会话授权，不再需要验证码。
	defer func() {
		_ = userApi.clearEmailCache(context, iRedis, LoginWithEmailPrefix, req.Email)
	}()

	users, err := s.emailCandidates(ctx, cache, req.Email, req.Code)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	options, _ := tenantPublicService.ListEnabledByIDs(util.DB(ctx), userService.TenantIdsOf(users))
	if len(users) > 1 {
		sessionId, serr := s.createEmailLoginSession(userService.UserIdsOf(users), options)
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

	token, err := s.buildToken(ctx, &users[0], options)
	if err != nil {
		response.ServerError(ctx)
		return
	}
	s.respondLoginResult(ctx, token, options)
}

// SendLoginWithEmailCode 发送使用邮箱登录验证码
func (s *AuthApi) SendLoginWithEmailCode(ctx *gin.Context) {
	req := &request.Email{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	users, err := userService.FindUsersByEmail(util.DB(ctx), req.Email)
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
	}, userService.UserIdsOf(users))
}
