package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HealthMonitoring struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Project        string         `gorm:"size:255;comment:项目名称" json:"project"`
	HttpStatus     int            `gorm:"comment:HTTP状态" json:"http_status"`
	Error          string         `gorm:"size:255;comment:错误信息" json:"error"`
	ResponseTimeMs int64          `gorm:"comment:响应时间(毫秒)" json:"response_time_ms"`
	CreatedAtUnix  int64          `gorm:"index;comment:创建时间(Unix秒)" json:"created_at"`
	UpdatedAtUnix  int64          `gorm:"comment:更新时间(Unix秒)" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// GORM Hook：在 Create 之前自动调用
func (h *HealthMonitoring) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().Unix()
	h.CreatedAtUnix = now
	h.UpdatedAtUnix = now
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

func (h *HealthMonitoring) BeforeUpdate(tx *gorm.DB) error {
	h.UpdatedAtUnix = time.Now().Unix()
	return nil
}
