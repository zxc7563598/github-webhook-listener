package health

import (
	"github.com/zxc7563598/github-webhook-listener/internal/config"
	service "github.com/zxc7563598/github-webhook-listener/internal/service/health"
)

// HealthHandler struct + 构造函数
type Handler struct {
	svc *service.Service
	cfg *config.Config
}

func New(svc *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		svc: svc,
		cfg: cfg,
	}
}
