package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	// 获取验证器自定义错误信息
	GetMessage() ValidatorMessages
}

// 验证器自定义错误信息字典
type ValidatorMessages map[string]string

// 获取自定义错误信息
func GetValidMsg(request Validator, err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, v := range validationErrors {
			key := v.Field() + "." + v.Tag()
			if message, exist := request.GetMessage()[key]; exist {
				return message
			}
			messages = append(messages, v.Error())
		}
		if len(messages) > 0 {
			return strings.Join(messages, "; ")
		}
	}

	return err.Error()
}
