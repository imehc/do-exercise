package cmd

import (
	"github.com/imehc/do-exercise/server/core"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/internal"
	"github.com/spf13/cobra"
)

var configFile string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动服务器",
	Long:  `启动do-exercise后端服务器。`,
	Run: func(cmd *cobra.Command, args []string) {
		initComponents()

		// 启动服务器
		defer shared.RSACrypto.Stop()
		defer shared.JobSchedulerInstance.Stop() // 停止定时任务调度器
		search := global.Searcher
		if search != nil {
			defer search.Close()
		}

		core.RunServer()
	},
}

func init() {
	// 添加配置文件标志
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "配置文件路径")

	// 添加server子命令
	rootCmd.AddCommand(serverCmd)
}

// initComponents 初始化所有服务
func initComponents() {
	internal.InitConfig(configFile)
	internal.InitLogger()
	internal.InitGorm(false)
	internal.InitRedis()
	internal.InitI18n()
	internal.InitCasbin()
	internal.InitMinio()
	internal.InitOther()
}
