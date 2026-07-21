package webhook

import (
	"fmt"
	"log"
	"time"

	"github.com/zxc7563598/github-webhook-listener/internal/config"
	"github.com/zxc7563598/github-webhook-listener/internal/queue"
)

func (s Service) MakeWebhookService(repoCfg *config.RepoConfig, branch string, event string, repoName string) error {
	log.Printf("[webhook] 仓库: %s, 事件: %s, 分支: %s", repoName, event, branch)
	for _, rule := range repoCfg.Rules {
		if rule.Event != event {
			continue
		}
		// 分支匹配
		if len(rule.Branches) > 0 {
			match := false
			for _, b := range rule.Branches {
				if b == branch {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		log.Printf("[webhook] 仓库 %s 的规则匹配: event=%s, branch=%s", repoName, rule.Event, branch)

		for _, action := range rule.Actions {
			if action.Type != "shell" {
				log.Printf("[webhook] 不支持的 action 类型: %s", action.Type)
				continue
			}
			taskName := fmt.Sprintf("%s:%s", repoName, rule.Event)
			taskID, err := s.repo.WebhookLogCreate(
				taskName,
				action.Command,
				action.WorkDir,
				action.Timeout,
				action.Env,
			)
			if err != nil {
				log.Printf("[database] 创建数据失败: %v", err)
				return fmt.Errorf("创建任务日志失败: %w", err)
			}
			if s.dispatcher == nil {
				log.Printf("[webhook] dispatcher 未初始化，无法调度任务")
				return fmt.Errorf("dispatcher 未初始化")
			}
			s.dispatcher.AddTask(&queue.ShellTask{
				ID:         taskID,
				Name:       taskName,
				Cmd:        action.Command,
				Timeout:    time.Duration(action.Timeout) * time.Second,
				RetryCount: action.RetryCount,
				RetryDelay: time.Duration(action.RetryDelay) * time.Second,
				Env:        action.Env,
				WorkDir:    action.WorkDir,
			})
		}
	}
	return nil
}
