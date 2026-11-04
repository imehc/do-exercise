package internal

import (
	"fmt"
	"path/filepath"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/util"
)

// InitOther 初始化其他工具
func InitOther() {
	initEmail()
	initJobScheduler()
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
