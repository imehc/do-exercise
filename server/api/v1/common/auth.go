package common

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/util"
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

	rsaCrypto := util.NewRSACrypto(global.Redis)
	Password, err := rsaCrypto.VerifyAndDecrypt(req.PublicKey, req.Password)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: err.Error(),
		})
		return
	}

	userId, err := userService.Login(common.Login{
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
		UserId:            userId,
		ExpireTime:        accessExpire,
		RefreshExpireTime: refreshExpire,
	}
	token, err := baseConf.GenerateToken()
	if err != nil {
		response.ServerError(ctx)
	}
	response.Success(ctx, token)
}

// PublicKey 获取公钥
func (s *AuthApi) PublicKey(ctx *gin.Context) {
	rsaCrypto := util.NewRSACrypto(global.Redis)
	publicKey, error := rsaCrypto.GenerateKeyPair()
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
