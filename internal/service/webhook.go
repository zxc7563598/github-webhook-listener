package service

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/zxc7563598/github-webhook-listener/internal/config"
	"github.com/zxc7563598/github-webhook-listener/internal/queue"
)

func (c Container) MakeWebhookService(repoCfg *config.RepoConfig, branch string, event string, repoName string) error {
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
		c.shellQueue.AddTask(&queue.ShellTask{
			ID:         uuid.NewString(),
			Name:       fmt.Sprintf("%s:%s", repoName, rule.Event),
			Cmd:        rule.Actions[0].Command,
			Args:       rule.Actions[0].Args,
			Timeout:    time.Duration(rule.Actions[0].Timeout) * time.Second,
			RetryCount: rule.Actions[0].RetryCount,
			RetryDelay: time.Duration(rule.Actions[0].RetryDelay) * time.Second,
			Env:        rule.Actions[0].Env,
			WorkDir:    rule.Actions[0].WorkDir,
		})
	}
	return nil
}
