package main

import (
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/internal"
	"github.com/imehc/do-exercise/server/middleware"
	"github.com/labstack/echo/v4"
)

func init() {
	internal.InitConfig()
	internal.InitLogger()
	// internal.InitGorm()
	// internal.InitRedis()
	internal.InitI18n()
}

func main() {

	// 创建Echo实例
	e := echo.New()

	// 添加日志中间件
	e.Use(middleware.Logger)

	// 添加i18n示例路由
	e.GET("/", func(c echo.Context) error {
		t := global.I18
		lang := c.FormValue("lang")
		// accept := c.Request().Header.Get("Accept-Language")
		// lang := c.Param("lang")
		// 使用指定语言翻译
		helloMsg := t.Translate("hello", "zh")
		welcomeMsg := t.Translate("welcome", lang)
		// 使用带变量的翻译
		greetingMsg := t.TranslateWithData("greeting", lang, map[string]interface{}{
			"Name": "Tom",
		})
		profileMsg := t.TranslateWithData("profile", lang, map[string]interface{}{
			"Name": "Tom",
			"Age":  25,
		})
		return c.JSON(200, map[string]string{
			"lang":     lang,
			"hello":    helloMsg,
			"welcome":  welcomeMsg,
			"greeting": greetingMsg,
			"profile":  profileMsg,
		})
	})

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", global.Config.System.Port)
	e.Logger.Fatal(e.Start(serverAddr))
}
