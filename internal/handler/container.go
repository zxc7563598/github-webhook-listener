package handler

import (
	"github.com/zxc7563598/github-webhook-listener/internal/config"
	"github.com/zxc7563598/github-webhook-listener/internal/service"
)

type Container struct {
	webEnabled        bool
	cfg               config.Config
	shellQueueService *service.Container
}

func NewContainer(webEnabled bool, cfg config.Config, shellQueueService *service.Container) *Container {
	return &Container{
		webEnabled:        webEnabled,
		cfg:               cfg,
		shellQueueService: shellQueueService,
	}
}
