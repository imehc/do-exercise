package internal

import (
	"fmt"
	"os"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/util"
	"github.com/spf13/viper"
)

// InitConfig 初始化配置
func InitConfig(configFile string) {
	// 创建viper实例
	v := viper.New()

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		util.Exit("配置文件不存在: ", err)
	}

	// 设置配置文件
	v.SetConfigFile(configFile)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到
			util.Exit("配置文件未找到: ", err)
		} else {
			// 配置文件存在但产生了其他错误
			util.Exit("配置文件错误: ", err)
		}
	}

	// 将配置映射到全局变量
	if err := v.Unmarshal(&global.Config); err != nil {
		util.Exit("配置映射失败: ", err)
	}
	fmt.Println("配置初始化完成")
}
