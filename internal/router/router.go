package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zxc7563598/github-webhook-listener/internal/config"
	"github.com/zxc7563598/github-webhook-listener/internal/handler"
	"github.com/zxc7563598/github-webhook-listener/internal/queue"
	"github.com/zxc7563598/github-webhook-listener/internal/service"
)

func SetupRouter(isWeb bool, cfg config.Config, scheduler *queue.ShellScheduler) *gin.Engine {
	r := gin.Default()
	// 中间件注册
	// r.Use()
	shellQueueService := service.NewContainer(scheduler)
	container := handler.NewContainer(isWeb, cfg, shellQueueService)

	r.GET("/", container.Lead)
	r.POST("/webhook", container.MakeWebhookHandler)
	if isWeb {
	}

	return r
}
