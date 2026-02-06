package webhook

import (
	"fmt"
	"time"

	webhookDTO "github.com/zxc7563598/github-webhook-listener/internal/dto/webhook"
	"github.com/zxc7563598/github-webhook-listener/pkg/utils"
)

func (s Service) QueryWebhookLogListPage(page, pageSize int, id *string) (*webhookDTO.WebhookLogListPage, error) {
	rows, total, err := s.repo.QueryWebhookLogListPage(page, pageSize, id)
	if err != nil {
		return nil, fmt.Errorf("数据获取失败: %v", err)
	}
	// 处理数据
	data := make([]webhookDTO.WebhookLogList, 0, len(rows))
	for _, r := range rows {
		data = append(data, webhookDTO.WebhookLogList{
			ID:        r.ID,
			Project:   r.Project,
			Command:   r.Command,
			Status:    r.Status,
			CreatedAt: time.Unix(r.CreatedAtUnix, 0).Format(time.DateTime),
		})
	}
	// 返回数据
	return &webhookDTO.WebhookLogListPage{
		Total: total,
		Data:  data,
	}, nil
}

func (s Service) QueryWebhookLogDetails(id string) (*webhookDTO.WebhookLogDetails, error) {
	row, err := s.repo.QueryWebhookLogDetails(id)
	if err != nil {
		return nil, fmt.Errorf("数据获取失败: %v", err)
	}
	// 处理数据
	startTime := ""
	endTime := ""
	exitCode := 0
	errString := ""
	stdout := ""
	stderr := ""
	if row.StartTimeUnix != nil {
		startTime = time.Unix(*row.StartTimeUnix, 0).Format(time.DateTime)
	}
	if row.EndTimeUnix != nil {
		endTime = time.Unix(*row.EndTimeUnix, 0).Format(time.DateTime)
	}
	if row.EndTimeUnix != nil {
		exitCode = *row.ExitCode
	}
	if row.Error != nil {
		errString = *row.Error
	}
	if row.StdoutPath != nil {
		stdout, err = utils.ReadLogFile(*row.StdoutPath)
		if err != nil {
		}
	}
	if row.StderrPath != nil {
		stderr, err = utils.ReadLogFile(*row.StderrPath)
		if err != nil {
		}
	}
	return &webhookDTO.WebhookLogDetails{
		ID:           row.ID,
		Project:      row.Project,
		Command:      row.Command,
		Status:       row.Status,
		StatusString: row.Status.String(),
		StartTime:    startTime,
		EndTime:      endTime,
		ExitCode:     exitCode,
		Error:        errString,
		Stdout:       stdout,
		Stderr:       stderr,
	}, nil

}
