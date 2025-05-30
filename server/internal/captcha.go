package internal

import (
	"image/color"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/util"
	"github.com/mojocn/base64Captcha"
)

// InitCaptcha 初始化验证码
func InitCaptcha() {
	// global.Captcha = base64Captcha.Se
	store := util.NewCaptcha(
		global.Redis,
		time.Second*time.Duration(global.Config.Captcha.CaptchaTimeOut),
	)
	captcheConf := global.Config.Captcha
	// 数字
	// driver := base64Captcha.NewDriverDigit(
	// 	captcheConf.ImgHeight,
	// 	captcheConf.ImgWidth,
	// 	captcheConf.KeyLong,
	// 	0.7,
	// 	80,
	// )
	// 文字
	driver := base64Captcha.NewDriverString(
		captcheConf.ImgHeight,
		captcheConf.ImgWidth,
		0,
		2|4,
		4,
		"1234567890qwertyuioplkjhgfdsazxcvbnm",
		&color.RGBA{R: 3, G: 102, B: 214, A: 125},
		base64Captcha.DefaultEmbeddedFonts,
		[]string{"wqy-microhei.ttc"},
	)
	// driver := base64Captcha.NewDriverMath(
	// 	captcheConf.ImgHeight,
	// 	captcheConf.ImgWidth,
	// 	0,
	// 	base64Captcha.OptionShowHollowLine, // 添加正弦线条增加辨识度
	// 	&color.RGBA{R: 145, G: 234, B: 145, A: 12}, // 深色背景
	// 	base64Captcha.DefaultEmbeddedFonts,
	// 	[]string{},
	// )
	global.Captcha = base64Captcha.NewCaptcha(driver, store)
}
