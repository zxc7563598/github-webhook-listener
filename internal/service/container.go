package service

import "github.com/zxc7563598/github-webhook-listener/internal/queue"

type Container struct {
	shellQueue *queue.ShellScheduler
}

func NewContainer(shellQueue *queue.ShellScheduler) *Container {
	return &Container{
		shellQueue: shellQueue,
	}
}
