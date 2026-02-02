package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zxc7563598/github-webhook-listener/internal/config"
	"github.com/zxc7563598/github-webhook-listener/internal/model"
)

// 创建 webhook 处理队列
func WebhookLogCreate(name, cmd, workDir string, timeout int, env []string) (string, error) {
	taskID := uuid.New()
	// 创建数据
	err := config.DB.Create(&model.WebhookLog{
		TaskID:     taskID,
		Project:    name,
		Command:    cmd,
		Env:        strings.Join(env, ";"),
		Timeout:    timeout,
		WorkDir:    workDir,
		Status:     model.WebhookLogStatusPending,
		RetryCount: 0,
	}).Error
	if err != nil {
		return taskID.String(), fmt.Errorf("创建webhook处理记录失败: %v", err)
	}
	return taskID.String(), nil
}

// webhook 处理失败重试
func WebhookLogRetryUpdate(taskID string, retryCount int) error {
	err := config.DB.Model(model.WebhookLog{}).Where("task_id = ?", taskID).Update("retry_count", retryCount).Error
	if err != nil {
		return fmt.Errorf("变更重试次数异常: %v", err)
	}
	return nil
}

// webhook 完成记录
func WebhookLogComplete(taskID, errMessage, stdout, stderr string, errCode int, status model.WebhookLogStatus, startTime, endTime time.Time) error {
	err := config.DB.Model(model.WebhookLog{}).Where("task_id = ?", taskID).Updates(map[string]any{
		"status":      status,
		"start_time":  startTime,
		"end_time":    endTime,
		"exit_code":   errCode,
		"error":       errMessage,
		"stdout_path": stdout,
		"stderr_path": stderr,
	}).Error
	if err != nil {
		return fmt.Errorf("存储完成信息异常: %v", err)
	}
	return nil
}
