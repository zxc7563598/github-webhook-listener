package repository

import (
	"fmt"

	"github.com/zxc7563598/github-webhook-listener/internal/config"
	"github.com/zxc7563598/github-webhook-listener/internal/model"
)

func HealthMonitoringCreate(project, errMessage string, httpStatus int, responseTimeMs int64) error {
	err := config.DB.Create(&model.HealthMonitoring{
		Project:        project,
		Error:          errMessage,
		HttpStatus:     httpStatus,
		ResponseTimeMs: responseTimeMs,
	}).Error
	if err != nil {
		return fmt.Errorf("健康监控信息记录失败: %v", err)
	}
	return nil
}
