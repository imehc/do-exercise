package config

// I18n 国际化配置
type I18n struct {
	DefaultLanguage string `yaml:"default_language" mapstructure:"default_language"` // 默认语言
	Languages      []string `yaml:"languages" mapstructure:"languages"`           // 支持的语言列表
}