package util

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/imehc/do-exercise/server/global"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

type CustomValidator struct {
	once      sync.Once
	validator *validator.Validate
}

func (c *CustomValidator) Validate(i any) error {
	c.lazyInit()
	if err := c.validator.Struct(i); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)

		// 获取请求的语言标识
		lang := language.English // 默认使用英语
		if ctx, ok := i.(echo.Context); ok {
			if langStr := ctx.Get("lang"); langStr != nil {
				lang = language.Make(langStr.(string))
			}
		}

		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			// 根据验证标签获取对应的错误消息
			message := global.I18.TranslateWithData(fmt.Sprintf("validation_%s", tag), lang.String(), map[string]interface{}{
				"Field": field,
			})

			if message == "" {
				message = global.I18.TranslateWithData("invalid_field", lang.String(), map[string]interface{}{
					"Field": field,
				})
			}

			errorMessages[field] = message
		}
		return echo.NewHTTPError(http.StatusBadRequest, errorMessages)
	}
	return nil
}

func (c *CustomValidator) lazyInit() {
	c.once.Do(func() {
		c.validator = validator.New()
		// 注册validate标签作为默认的验证标签
		c.validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("validate"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	})
}
