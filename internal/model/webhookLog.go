package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebhookLogStatus int

const (
	WebhookLogStatusPending WebhookLogStatus = iota
	WebhookLogStatusRunning
	WebhookLogStatusSuccess
	WebhookLogStatusFailed
	WebhookLogStatusTimeout
	WebhookLogStatusCancelled
)

func (status WebhookLogStatus) String() string {
	switch status {
	case WebhookLogStatusPending:
		return "待处理"
	case WebhookLogStatusRunning:
		return "运行中"
	case WebhookLogStatusSuccess:
		return "成功"
	case WebhookLogStatusFailed:
		return "失败"
	case WebhookLogStatusTimeout:
		return "超时"
	case WebhookLogStatusCancelled:
		return "取消"
	default:
		return "未知"
	}
}

type WebhookLog struct {
	ID            uuid.UUID        `gorm:"type:char(36);primaryKey" json:"id"`
	TaskID        uuid.UUID        `gorm:"type:char(36)" json:"task_id"`
	Project       string           `gorm:"size:255;comment:项目名称" json:"project"`
	Command       string           `gorm:"size:255;comment:执行命令" json:"command"`
	Env           string           `gorm:"size:255;comment:环境变量" json:"env"`
	Timeout       int              `gorm:"comment:允许执行时间，默认 300 单位秒" json:"timeout"`
	WorkDir       string           `gorm:"size:255;comment:命令工作目录，默认项目二进制文件目录" json:"work_dir"`
	Status        WebhookLogStatus `gorm:"comment:任务状态" json:"status"`
	StartTimeUnix *int64           `gorm:"comment:开始时间" json:"start_time_unix"`
	EndTimeUnix   *int64           `gorm:"comment:结束时间" json:"end_time_unix"`
	ExitCode      *int             `gorm:"comment:退出code" json:"exit_code"`
	Error         *string          `gorm:"size:255;comment:错误信息" json:"error"`
	RetryCount    int              `gorm:"comment:重试次数" json:"retry_count"`
	StdoutPath    *string          `gorm:"size:255;comment:Stdout 日志路径" json:"stdout_path"`
	StderrPath    *string          `gorm:"size:255;comment:Stderr 日志路径" json:"stderr_path"`
	CreatedAtUnix int64            `gorm:"index;comment:创建时间(Unix秒)" json:"created_at_unix"`
	UpdatedAtUnix int64            `gorm:"comment:更新时间(Unix秒)" json:"updated_at_unix"`
	DeletedAt     gorm.DeletedAt   `gorm:"index" json:"deleted_at"`
}

// GORM Hook：在 Create 之前自动调用
func (w *WebhookLog) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().Unix()
	w.CreatedAtUnix = now
	w.UpdatedAtUnix = now
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

func (w *WebhookLog) BeforeUpdate(tx *gorm.DB) error {
	w.UpdatedAtUnix = time.Now().Unix()
	return nil
}
