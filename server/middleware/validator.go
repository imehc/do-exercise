package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/internal"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
)

// ValidatorMiddleware 验证错误处理中间件
func ValidaterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取语言设置，默认为英语
		lang := getLanguage(c.GetHeader("Accept-Language"))
		if lang == "" || (lang != "en" && lang != "zh") {
			lang = "en"
		}
		c.Set("lang", lang)
		trans, err := internal.InitTrans(lang)
		if err != nil {
			response.ServerError(c)
		}
		c.Next()

		errors := c.Errors.Last()
		if errors != nil {
			if validationErrors, ok := errors.Err.(validator.ValidationErrors); ok {
				translatedErrors := validationErrors.Translate(trans)
				details := make([]response.ValidationDetail, 0)
				for field, message := range removeTopStruct(translatedErrors) {
					details = append(details, response.ValidationDetail{
						Field:    field,
						Messages: []string{message},
					})
				}
				response.BadRequest(c, response.ValidationError{
					Type:    status.BAD_REQUEST_MSG,
					Message: global.I18.Translate("badRequest", lang),
					Details: details,
				})
				return
			}
			response.ServerError(c)
			return
		}
	}
}

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func getLanguage(originLang string) string {
	lang := strings.ToLower(originLang)

	switch {
	case strings.HasPrefix(lang, "zh"):
		return "zh"
	case strings.HasPrefix(lang, "en"):
		return "en"
	default:
		return "en" // default fallback
	}
}
