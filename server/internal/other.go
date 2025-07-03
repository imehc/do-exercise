package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/util"
)

// InitOther 初始化其他工具
func InitOther() {
	initSnowflake()
	initEmail()
	initJobScheduler()
}

func initSnowflake() {
	nodeID := int64(1) // 默认值
	if v := os.Getenv("SNOWFLAKE_NODE_ID"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			nodeID = id
		}
	}
	if err := util.InitSnowflake(nodeID); err != nil {
		util.Exit("初始化雪花算法失败: ", err)
	}
}

func initEmail() {
	shared.Email = &util.Email{
		TemplatePath: filepath.Join("template", "email.html"),
		SmtpHost:     global.Config.Email.Host,
		SmtpPort:     global.Config.Email.Port,
		SmtpUser:     global.Config.Email.User,
		SmtpPass:     global.Config.Email.Pass,
	}
}

func initJobScheduler() {
	shared.JobSchedulerInstance = shared.GetJobScheduler()

	// 恢复数据库中的定时任务
	if err := shared.JobSchedulerInstance.RestoreJobsFromDatabase(); err != nil {
		fmt.Printf("恢复定时任务失败: %v\n", err)
	} else {
		fmt.Println("定时任务恢复完成")
	}
}
