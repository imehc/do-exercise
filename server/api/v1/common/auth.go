package common

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/util"
)

type AuthApi struct{}

// Login 登录
func (s *AuthApi) Login(ctx *gin.Context) {
	var req common.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if ok, err := userService.Login(common.Login{}); err != nil || !ok {
		lang := ctx.GetString("lang")
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
		UserId:            111,
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
