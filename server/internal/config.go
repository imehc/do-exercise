package internal

import (
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/spf13/viper"
)

// InitConfig 初始化配置
func InitConfig() {
	// 创建viper实例
	v := viper.New()

	// 设置配置文件名
	v.SetConfigName("config")
	// 设置配置文件类型
	v.SetConfigType("yaml")
	// 设置配置文件路径
	v.AddConfigPath(".")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到
			panic(fmt.Errorf("配置文件未找到: %v", err))
		} else {
			// 配置文件存在但产生了其他错误
			panic(fmt.Errorf("读取配置文件失败: %v", err))
		}
	}

	// 将配置映射到全局变量
	if err := v.Unmarshal(&global.Config); err != nil {
		panic(fmt.Errorf("配置映射失败: %v", err))
	}
	fmt.Println("配置初始化完成")
}
