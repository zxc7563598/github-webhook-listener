package webhook

import (
	"strings"

	"github.com/google/uuid"
	"github.com/zxc7563598/github-webhook-listener/internal/model"
	"gorm.io/gorm"
)

type gormRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) WebhookLogRepository {
	return &gormRepo{db: db}
}

// 创建 webhook 处理队列
func (r *gormRepo) WebhookLogCreate(name, cmd, workDir string, timeout int, env []string) (string, error) {
	taskID := uuid.New()
	err := r.db.Create(&model.WebhookLog{
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
		return taskID.String(), err
	}
	return taskID.String(), nil
}

// webhook 处理失败重试
func (r *gormRepo) WebhookLogRetryUpdate(taskID string, retryCount int) error {
	err := r.db.Model(&model.WebhookLog{}).Where("task_id = ?", taskID).Update("retry_count", retryCount).Error
	if err != nil {
		return err
	}
	return nil
}

// webhook 完成记录
func (r *gormRepo) WebhookLogComplete(taskID, errMessage, stdout, stderr string, errCode int, status model.WebhookLogStatus, startTime, endTime int64) error {
	err := r.db.Model(&model.WebhookLog{}).Where("task_id = ?", taskID).Updates(map[string]any{
		"status":          status,
		"start_time_unix": startTime,
		"end_time_unix":   endTime,
		"exit_code":       errCode,
		"error":           errMessage,
		"stdout_path":     stdout,
		"stderr_path":     stderr,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// 查询分页列表
type webhookLogListRow struct {
	ID            string                 `gorm:"column:id"`
	Project       string                 `gorm:"column:project"`
	Command       string                 `gorm:"column:command"`
	Status        model.WebhookLogStatus `gorm:"column:status"`
	CreatedAtUnix int64                  `gorm:"column:created_at_unix"`
}

func (r *gormRepo) QueryWebhookLogListPage(page, pageSize int, id *string) ([]webhookLogListRow, int64, error) {
	var rows []webhookLogListRow
	var total int64
	var err error
	limit := pageSize
	offset := (page - 1) * pageSize
	baseDB := r.db.Model(&model.WebhookLog{})
	if id != nil && *id != "" {
		baseDB = baseDB.Where("id = ?", *id)
	}
	listDB := baseDB
	countDB := baseDB
	err = countDB.Count(&total).Error
	if err != nil {
		return rows, total, err
	}
	if total == 0 {
		return rows, total, nil
	}
	err = listDB.Select("id, project, command, status, created_at_unix").Order("created_at_unix DESC").Limit(limit).Offset(offset).Scan(&rows).Error
	if err != nil {
		return rows, total, err
	}
	return rows, total, nil
}

type webhookLogDetailsRow struct {
	ID            string                 `gorm:"column:id"`
	Project       string                 `gorm:"column:project"`
	Command       string                 `gorm:"column:command"`
	Status        model.WebhookLogStatus `gorm:"column:status"`
	StartTimeUnix *int64                 `gorm:"column:start_time_unix"`
	EndTimeUnix   *int64                 `gorm:"column:end_time_unix"`
	ExitCode      *int                   `gorm:"column:exit_code"`
	Error         *string                `gorm:"column:error"`
	StdoutPath    *string                `gorm:"column:stdout_path"`
	StderrPath    *string                `gorm:"column:stderr_path"`
}

func (r *gormRepo) QueryWebhookLogDetails(id string) (*webhookLogDetailsRow, error) {
	var row webhookLogDetailsRow
	tx := r.db.Model(&model.WebhookLog{}).Where("id = ?", id).Select("id, project, command, status, start_time_unix, end_time_unix, exit_code, error, stdout_path, stderr_path").Scan(&row)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}
