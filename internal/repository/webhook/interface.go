package webhook

import (
	"github.com/zxc7563598/github-webhook-listener/internal/model"
)

// WebhookLogRepository 定义 webhook 日志的数据访问接口
type WebhookLogRepository interface {
	WebhookLogCreate(name, cmd, workDir string, timeout int, env []string) (string, error)
	WebhookLogRetryUpdate(taskID string, retryCount int) error
	WebhookLogComplete(taskID, errMessage, stdout, stderr string, errCode int, status model.WebhookLogStatus, startTime, endTime int64) error
	QueryWebhookLogListPage(page, pageSize int, id *string) ([]webhookLogListRow, int64, error)
	QueryWebhookLogDetails(id string) (*webhookLogDetailsRow, error)
}
