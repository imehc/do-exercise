package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/imehc/do-exercise/server/global"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// InitI18n 初始化国际化翻译器
func InitI18n() {
	// 创建一个新的Bundle，设置默认语言为英语
	bundle := i18n.NewBundle(language.English)
	// 注册yaml解析器
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// 获取i18n目录
	i18nDir := "i18n"
	// 加载所有yaml文件
	err := filepath.Walk(i18nDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			_, err = bundle.LoadMessageFile(path)
			if err != nil {
				panic(fmt.Errorf("加载国际化文件失败 %s: %v", path, err))
			}
		}
		return nil
	})

	if err != nil {
		panic(fmt.Errorf("遍历国际化目录失败: %v", err))
	}

	// 初始化 global.I18n
	localizer := i18n.NewLocalizer(bundle, language.English.String())
	global.I18 = &global.I18n{
		Bundle:    bundle,
		Localizer: localizer,
	}

	fmt.Println("国际化初始化完成")
}
