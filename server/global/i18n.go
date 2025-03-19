package global

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type I18n struct {
	*i18n.Bundle
	*i18n.Localizer
}

// Translate 翻译消息
func (i I18n) Translate(messageID string, lang string) string {
	// 为当前请求创建一个新的 localizer
	localizer := i18n.NewLocalizer(i.Bundle, lang)

	// 翻译消息
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID // 如果翻译失败，返回原始消息ID
	}

	return msg
}

// TranslateWithData 翻译带变量的消息
func (i I18n) TranslateWithData(messageID string, lang string, templateData map[string]interface{}) string {
	// 为当前请求创建一个新的 localizer
	localizer := i18n.NewLocalizer(i.Bundle, lang)

	// 翻译消息
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		return messageID // 如果翻译失败，返回原始消息ID
	}

	return msg
}
