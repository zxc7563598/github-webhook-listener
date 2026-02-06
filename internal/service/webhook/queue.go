package webhook

import (
	"errors"

	"github.com/zxc7563598/github-webhook-listener/internal/model"
	"github.com/zxc7563598/github-webhook-listener/internal/queue"
)

func (s *Service) SetDispatcher(d queue.ShellTaskDispatcher) {
	s.dispatcher = d
}

func (s *Service) WebhookLogRetryUpdate(taskID string, retryCount int) error {
	err := s.repo.WebhookLogRetryUpdate(taskID, retryCount)
	return err
}

func (s *Service) WebhookLogComplete(taskID, errMessage, stdout, stderr string, errCode int, status model.WebhookLogStatus, startTime, endTime int64) error {
	err := s.repo.WebhookLogComplete(taskID, errMessage, stdout, stderr, errCode, status, startTime, endTime)
	return err
}

func (s *Service) Trigger(task *queue.ShellTask) error {
	if s.dispatcher == nil {
		return errors.New("dispatcher not initialized")
	}
	return s.dispatcher.AddTask(task)
}
