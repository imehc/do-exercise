package internal

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var lock sync.Mutex // 添加互斥锁

// 定义翻译的方法
func InitTrans(locale string) (trans ut.Translator, err error) {
	lock.Lock()
	defer lock.Unlock()
	//修改gin框架中的validator引擎属性, 实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		//第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			panic(fmt.Errorf("uni.GetTranslator(%s)", locale))
		}
		// 注册自定义校验函数
		v.RegisterValidation("startWithLetter", startWithLetter) // 校验以字母开头
		v.RegisterValidation("containsLetter", containsLetter)   // 校验至少包含一个字母
		v.RegisterValidation("complexPassword", complexPassword) // 校验密码是否包含字母、数字和特殊字符
		v.RegisterValidation("uniqueName", uniqueName)           // 校验可以为空，但是如果不为空则必须为唯一，且只能包含字母、数字和冒号
		v.RegisterValidation("menuType", menuType)               // 校验菜单类型是否为1或2或3

		switch locale {
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, trans)
			v.RegisterTranslation("startWithLetter", trans, func(ut ut.Translator) error {
				return ut.Add("startWithLetter", "{0}必须以字母开头", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("startWithLetter", fe.Field())
				return t
			})
			v.RegisterTranslation("containsLetter", trans, func(ut ut.Translator) error {
				return ut.Add("containsLetter", "{0}必须包含至少一个字母", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("containsLetter", fe.Field())
				return t
			})
			v.RegisterTranslation("complexPassword", trans, func(ut ut.Translator) error {
				return ut.Add("complexPassword", "{0}必须包含字母、数字和特殊字符", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("complexPassword", fe.Field())
				return t
			})
			v.RegisterTranslation("uniqueName", trans, func(ut ut.Translator) error {
				return ut.Add("uniqueName", "{0}为空，或者必须为唯一，且只能包含字母、数字和冒号", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("uniqueName", fe.Field())
				return t
			})
			v.RegisterTranslation("menuType", trans, func(ut ut.Translator) error {
				return ut.Add("menuType", "{0}必须为1或2或3", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("menuType", fe.Field())
				return t
			})
		default:
			en_translations.RegisterDefaultTranslations(v, trans)
			v.RegisterTranslation("startWithLetter", trans, func(ut ut.Translator) error {
				return ut.Add("startWithLetter", "{0} must start with a letter", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("startWithLetter", fe.Field())
				return t
			})
			v.RegisterTranslation("containsLetter", trans, func(ut ut.Translator) error {
				return ut.Add("containsLetter", "{0} must contain at least one letter", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("containsLetter", fe.Field())
				return t
			})
			v.RegisterTranslation("complexPassword", trans, func(ut ut.Translator) error {
				return ut.Add("complexPassword", "{0} must contain letters, numbers, and special characters", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("complexPassword", fe.Field())
				return t
			})
			v.RegisterTranslation("uniqueName", trans, func(ut ut.Translator) error {
				return ut.Add("uniqueName", "{0} can be empty, or must be unique, and can only contain letters, numbers, and colons", true)

			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("uniqueName", fe.Field())
				return t
			})
			v.RegisterTranslation("menuType", trans, func(ut ut.Translator) error {
				return ut.Add("menuType", "{0} must be 1 or 2 or 3", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("menuType", fe.Field())
				return t
			})
		}
		return
	}
	return
}

// 自定义校验函数：检查字符串是否以字母开头
func startWithLetter(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) == 0 {
		return false
	}
	// 使用正则表达式检查第一个字符是否为字母
	match, _ := regexp.MatchString(`^[a-zA-Z]`, string(value[0]))
	return match
}

// 自定义校验函数：检查字符串是否至少包含一个字母
func containsLetter(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// 使用正则表达式检查是否包含字母
	match, _ := regexp.MatchString(`[a-zA-Z]`, value)
	return match
}

// 自定义校验函数：检查密码是否包含字母、数字和特殊字符
func complexPassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// 使用正则表达式检查密码是否包含字母、数字和特殊字符
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(value)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(value)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(value)

	return hasLetter && hasDigit && hasSpecial
}

// 自定义校验函数：检查可以为空，但是如果不为空则必须为唯一，且只能包含字母、数字和冒号
func uniqueName(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	// 使用正则表达式检查字符串是否只包含字母、数字和冒号
	match, _ := regexp.MatchString(`^[a-zA-Z0-9:]+$`, value)
	return match
}

// 自定义校验函数：检查菜单类型是否为1或2或3
func menuType(fl validator.FieldLevel) bool {
	value := fl.Field().Uint()
	return value == 1 || value == 2 || value == 3
}
