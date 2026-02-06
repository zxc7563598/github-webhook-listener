package webhook

import "github.com/zxc7563598/github-webhook-listener/internal/model"

type WebhookLogList struct {
	ID        string
	Project   string
	Command   string
	Status    model.WebhookLogStatus
	CreatedAt string
}

type WebhookLogListPage struct {
	Total int64
	Data  []WebhookLogList
}

type WebhookLogDetails struct {
	ID           string
	Project      string
	Command      string
	Status       model.WebhookLogStatus
	StatusString string
	StartTime    string
	EndTime      string
	ExitCode     int
	Error        string
	Stdout       string
	Stderr       string
}
