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
	ResponseTimeMs int64          `gorm:"comment:响应时间(毫秒)" json:"interval"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// GORM Hook：在 Create 之前自动调用
func (h *HealthMonitoring) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}
