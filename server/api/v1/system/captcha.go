package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	sysRes "github.com/imehc/do-exercise/server/model/system/response"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type CaptchaApi struct{}

// @Summary 获取验证码
// @Description 获取验证码
// @Tags 验证码
// @Accept json
// @Produce json
// @Success 200 {object} response.Captcha "验证码"
// @Failure 400 {object} string "验证码获取失败"
// @Router /captcha [get]
func (c *CaptchaApi) GetCaptcha(ctx *gin.Context) {
	imageHeight := global.CONFIG.Captcha.ImgHeight
	imageWidth := global.CONFIG.Captcha.ImgWidth
	keyLong := global.CONFIG.Captcha.KeyLong
	driver := base64Captcha.NewDriverDigit(imageHeight, imageWidth, keyLong, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, global.CAPTCHA_STORE)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		global.LOG.Error("验证码获取失败!", zap.Error(err))
		response.BadRequest("验证码获取失败", ctx)
		return
	}

	response.Success(sysRes.Captcha{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: keyLong,
	}, ctx)
}
