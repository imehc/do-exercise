package common

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type AuthApi struct{}

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

	Password, err := shared.RSACrypto.DecryptWithKey(req.PublicKey, req.Password)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return
	}

	user, err := authService.Login(common.Login{
		Username: req.Username,
		Password: Password,
	})
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate("usernameOrPasswordError", lang),
		})
		return
	}
	accessExpire, err := util.ParseDurationString(global.Config.Auth.AccessExpireTime)
	if err != nil {
		response.ServerError(ctx)
	}
	refreshExpire, err := util.ParseDurationString(global.Config.Auth.RefreshExpireTime)
	if err != nil {
		response.ServerError(ctx)
	}
	baseConf := util.Token{
		UserId: user.Id,
		RoleIds: lo.Map(user.Roles, func(item system.SysRole, index int) uint {
			return item.Id
		}),
		ExpireTime:        accessExpire,
		RefreshExpireTime: refreshExpire,
	}
	token, err := baseConf.GenerateToken()
	if err != nil {
		response.ServerError(ctx)
	}
	response.Success(ctx, token)
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
