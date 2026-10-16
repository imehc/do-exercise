package system

import (
	"time"

	"github.com/imehc/do-exercise/server/model"
	"gorm.io/gorm"
)

// SysJob 定时任务
type SysJob struct {
	model.IdWrapper

	Name           string     `json:"name" gorm:"size:128;not null;comment:任务名称"`
	JobGroup       string     `json:"job_group" gorm:"size:128;not null;comment:任务分组"`
	CronExpression string     `json:"cron_expression" gorm:"size:64;not null;comment:cron执行表达式"`
	Command        string     `json:"command" gorm:"size:2048;not null;comment:执行命令"`
	Status         uint8      `json:"status" gorm:"default:1;comment:状态(1:正常,2:暂停)"`
	Concurrent     bool       `json:"concurrent" gorm:"default:false;comment:是否并发执行(true:是,false:否)"`
	Description    string     `json:"description" gorm:"size:512;comment:任务描述"`
	LastTime       *time.Time `json:"last_time" gorm:"comment:上次执行时间"`
	NextTime       *time.Time `json:"next_time" gorm:"comment:下次执行时间"`
	Times          int64      `json:"times" gorm:"default:0;comment:执行次数"`
	RetryTimes     uint       `json:"retry_times" gorm:"default:0;comment:重试次数"`
	RetryInterval  uint       `json:"retry_interval" gorm:"default:0;comment:重试间隔(秒)"`
	Timeout        uint       `json:"timeout" gorm:"default:0;comment:执行超时时间(秒)"`
	TenantId       string     `json:"tenant_id" gorm:"column:tenant_id;type:varchar(32);not null;default:'';index;comment:租户ID"`

	model.ControlWrapper
}

func (j *SysJob) BeforeCreate(tx *gorm.DB) (err error) {
	userId := model.CurrentUserID(tx)
	j.CreatedBy = userId
	j.UpdatedBy = userId

	return nil
}

func (j *SysJob) BeforeUpdate(tx *gorm.DB) (err error) {
	userId := model.CurrentUserID(tx)
	if userId == "" {
		// 没有用户ID，直接跳过
		return nil
	}

	if j.UpdatedBy != userId && j.Id != 0 {
		j.UpdatedBy = userId
		err = tx.
			Model(j).
			Select("UpdatedBy").
			Updates(j).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *SysJob) BeforeDelete(tx *gorm.DB) (err error) {
	if j.Id != 0 {
		userId := model.CurrentUserID(tx)
		j.DeletedBy = userId
		err = tx.
			Model(j).
			Select("DeletedBy").
			Updates(j).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
