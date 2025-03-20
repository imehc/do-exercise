package router

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Run() *gin.Engine {
	r := gin.Default()
	// 获取 validator 实例
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义校验函数
		v.RegisterValidation("startWithLetter", startWithLetter) // 校验以字母开头
		v.RegisterValidation("containsLetter", containsLetter)   // 校验至少包含一个字母
		v.RegisterValidation("complexPassword", complexPassword) // 校验密码是否包含字母、数字和特殊字符
	}

	r.Use(gin.Recovery())

	system := RouterGroupApp.System

	protected := r.Group("/system")
	public := r.Group("/")
	{
		// 健康监测
		public.GET("/health", func(c *gin.Context) {

		})
	}
	{
		system.InitSysUserRouter(protected)
	}

	return r
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
