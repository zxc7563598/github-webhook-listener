package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/github-webhook-listener/internal/config"
	"github.com/zxc7563598/github-webhook-listener/internal/handler"
	healthHandler "github.com/zxc7563598/github-webhook-listener/internal/handler/health"
	webhookHandler "github.com/zxc7563598/github-webhook-listener/internal/handler/webhook"
	"github.com/zxc7563598/github-webhook-listener/internal/queue"
	healthRepo "github.com/zxc7563598/github-webhook-listener/internal/repository/health"
	webhookRepo "github.com/zxc7563598/github-webhook-listener/internal/repository/webhook"
	healthSvc "github.com/zxc7563598/github-webhook-listener/internal/service/health"
	webhookSvc "github.com/zxc7563598/github-webhook-listener/internal/service/webhook"
)

func NewApp(isWeb *bool, cfg *config.Config, webUser, webPass *string) (*gin.Engine, *queue.ShellScheduler, *queue.HealthMonitor) {
	// repository
	webhookRepo := webhookRepo.New(config.DB)
	healthRepo := healthRepo.New(config.DB)

	// service
	webhookSvc := webhookSvc.New(webhookRepo)
	healthSvc := healthSvc.New(healthRepo, webhookSvc)

	// queue（注入 service 作为结果处理者）
	scheduler := queue.NewShellScheduler(5, webhookSvc)
	scheduler.Start()

	healthMonitor := queue.NewHealthMonitor(cfg, healthSvc)
	healthMonitor.Start()

	// 反向注入
	webhookSvc.SetDispatcher(scheduler)

	// handler
	webhookHandler := webhookHandler.New(webhookSvc, cfg)
	healthHandler := healthHandler.New(healthSvc, cfg)

	r := gin.New()
	handler.Register(isWeb, webUser, webPass, r, webhookHandler, healthHandler)

	return r, scheduler, healthMonitor
}
