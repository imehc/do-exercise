package initialize

import (
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/mojocn/base64Captcha"
)

// 初始化其他
func InitOther() {
	global.CAPTCHA_STORE = base64Captcha.NewMemoryStore(10240, time.Second*time.Duration(global.CONFIG.Captcha.CaptchaTimeOut))
}
