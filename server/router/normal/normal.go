package normal

import (
	"github.com/imehc/do-exercise/server/global"
	"github.com/labstack/echo/v4"
)

type NormalRouter struct {
}

func (n *NormalRouter) InitNormalRouter(r *echo.Group) {
	router := r.Group("")
	{
		router.GET("/", func(c echo.Context) error {
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
	}
}
