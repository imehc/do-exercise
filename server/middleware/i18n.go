package middleware

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

// I18nMiddleware 处理国际化中间件
func I18nMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 默认使用英语
			lang := language.English

			if acceptLang := c.Request().Header.Get("Accept-Language"); acceptLang != "" {
				// 从请求头获取Accept-Language
				if tags, _, err := language.ParseAcceptLanguage(acceptLang); err == nil && len(tags) > 0 {
					lang = tags[0]
				}
			}

			// 将语言设置存储在上下文中
			c.Set("lang", lang.String())

			return next(c)
		}
	}
}
