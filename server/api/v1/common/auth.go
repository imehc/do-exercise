package common

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type AuthApi struct{}

func (s *AuthApi) generateToken(ctx *gin.Context, user *system.SysUser) {
	accessExpire, err := util.ParseDurationString(global.Config.Auth.AccessExpireTime)
	if err != nil {
		response.ServerError(ctx)
	}
	refreshExpire, err := util.ParseDurationString(global.Config.Auth.RefreshExpireTime)
	if err != nil {
		response.ServerError(ctx)
	}
	baseConf := util.Token{
		UserId:   user.UserId,
		Username: user.Username,
		RoleIds: lo.Map(user.Roles, func(item system.SysRole, index int) uint {
			return item.Id
		}),
		ExpireTime:        accessExpire,
		RefreshExpireTime: refreshExpire,
		Disabled:          false, // TODO: 在数据库中添加字段获取禁用状态
		CreatedTime:       time.Now(),
	}
	token, err := baseConf.GenerateToken()
	if err != nil {
		response.ServerError(ctx)
	}
	response.Success(ctx, token)
}

// Login 登录
func (s *AuthApi) Login(ctx *gin.Context) {
	lang := ctx.GetString("lang")

	var req common.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if ok := global.Captcha.Verify(req.CaptchaId, req.Captcha, true); !ok {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("invalidCaptcha", lang),
		})
		return
	}

	password, err := shared.RSACrypto.DecryptWithKey(req.PublicKey, req.Password, true)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}

	user, err := authService.Login(common.Login{
		Username: req.Username,
		Password: password,
	})
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}

	s.generateToken(ctx, user)
}

// RefreshToken 刷新token
func (s *AuthApi) RefreshToken(ctx *gin.Context) {
	lang := ctx.GetString("lang")

	type RefreshTokenReq struct {
		RefreshToken string `form:"refresh_token" binding:"required"`
	}
	req := RefreshTokenReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("emptyRefreshToken", lang),
		})
		return
	}
	accessExpire, err := util.ParseDurationString(global.Config.Auth.AccessExpireTime)
	if err != nil {
		response.ServerError(ctx)
	}
	baseConf := util.Token{
		ExpireTime: accessExpire,
	}
	token, err := baseConf.RefreshToken(req.RefreshToken)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
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
	lang := ctx.GetString("lang")
	id, b64s, _, err := global.Captcha.Generate()
	if err != nil {
		global.Log.Error("验证码获取失败!", zap.Error(err))
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("captchaError", lang),
		})
		return
	}
	response.Success(ctx, common.Captcha{
		CaptchaId: id,
		PicPath:   b64s,
	})
}

// ResetPassword 忘记密码
func (s *AuthApi) ResetPassword(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	iRedis := global.Redis
	context := context.Background()

	req := &request.UserResetPasswordReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	user, err := userService.FindUserByEmail(req.Email)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("emailNotBound", lang),
		})
		return
	}

	password, err := shared.RSACrypto.DecryptWithKey(req.PublicKey, req.Password, true)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}

	var userApi UserApi
	cache, err := userApi.getEmailCache(context, iRedis, ForgotPasswordPrefix, req.Email)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}

	defer func() {
		_ = userApi.clearEmailCache(context, iRedis, ForgotPasswordPrefix, req.Email)
	}()

	if cache.Code != req.Code || cache.UserId != user.UserId {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("captchaError", lang),
		})
		return
	}

	if err := authService.ResetPassword(request.UserResetPasswordReq{
		Id:       user.UserId,
		Password: password,
	}); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}
	response.NoContent(ctx)
}

// SendResetPasswordCode 发送忘记密码邮箱验证码
func (s *AuthApi) SendResetPasswordCode(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	req := &request.Email{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
	}

	user, err := userService.FindUserByEmail(req.Email)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
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
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), ctx.GetString("lang")),
		})
		return
	}

	response.NoContent(ctx)
}

// LoginWithEmail 使用邮件登录
func (s *AuthApi) LoginWithEmail(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	context := context.Background()
	iRedis := global.Redis

	var req common.LoginWithEmailReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
	}

	var userApi UserApi
	cache, err := userApi.getEmailCache(context, iRedis, LoginWithEmailPrefix, req.Email)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}

	defer func() {
		_ = userApi.clearEmailCache(context, iRedis, LoginWithEmailPrefix, req.Email)
	}()

	user, err := userService.FindUserByEmail(req.Email)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("emailNotBound", lang),
		})
		return
	}

	if cache.Code != req.Code || cache.UserId != user.UserId {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("captchaError", lang),
		})
		return
	}

	s.generateToken(ctx, user)
}

// SendLoginWithEmailCode 发送使用邮箱登录验证码
func (s *AuthApi) SendLoginWithEmailCode(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	req := &request.Email{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
	}

	user, err := userService.FindUserByEmail(req.Email)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}

	var userApi UserApi
	userApi.sendEmailCode(ctx, LoginWithEmailPrefix, &util.EmailData{
		EmailTitle:       "验证码",
		VerificationType: "邮箱登录",
		GreetingText:     "感谢您使用我们的服务，请使用以下验证码完成登录：",
	}, user.UserId)
}
