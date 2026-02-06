package health

import (
	repository "github.com/zxc7563598/github-webhook-listener/internal/repository/health"
	webhookService "github.com/zxc7563598/github-webhook-listener/internal/service/webhook"
)

// HealthService struct + 构造函数
type Service struct {
	repo       repository.HealthRepository
	webhookSvc webhookService.ReadOnlyWebhookService
}

func New(repo repository.HealthRepository, webhookSvc webhookService.ReadOnlyWebhookService) *Service {
	return &Service{
		repo:       repo,
		webhookSvc: webhookSvc,
	}
}

func (s Service) HealthMonitoringCreate(project, errMessage string, httpStatus int, responseTimeMs int64) error {
	err := s.repo.HealthMonitoringCreate(project, errMessage, httpStatus, responseTimeMs)
	return err
}
