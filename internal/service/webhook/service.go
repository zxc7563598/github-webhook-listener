package webhook

import (
	webhookDTO "github.com/zxc7563598/github-webhook-listener/internal/dto/webhook"
	"github.com/zxc7563598/github-webhook-listener/internal/queue"
	repository "github.com/zxc7563598/github-webhook-listener/internal/repository/webhook"
)

// WebhookService struct + 构造函数
type Service struct {
	repo       repository.WebhookLogRepository
	dispatcher queue.ShellTaskDispatcher
}

func New(repo repository.WebhookLogRepository) *Service {
	return &Service{
		repo: repo,
	}
}

type ReadOnlyWebhookService interface {
	QueryWebhookLogListPage(page, pageSize int, id *string) (*webhookDTO.WebhookLogListPage, error)
	QueryWebhookLogDetails(id string) (*webhookDTO.WebhookLogDetails, error)
}
