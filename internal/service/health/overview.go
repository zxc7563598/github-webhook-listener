package health

import (
	"math"
	"net/http"

	"github.com/zxc7563598/github-webhook-listener/internal/config"
	healthDTO "github.com/zxc7563598/github-webhook-listener/internal/dto/health"
)

func (s Service) GetOverview(repos map[string]*config.RepoConfig) (healthDTO.Overview, error) {
	var result healthDTO.Overview
	var projects = make(map[string]healthDTO.OverviewProject)
	// 获取配置过的项目信息
	for repoName, repoCfg := range repos {
		result.ProjectCount++
		// 定义项目数据
		history := make([]float64, 24)
		for i := range history {
			history[i] = -1
		}
		projects[repoName] = healthDTO.OverviewProject{
			ID:               repoName,
			Name:             repoCfg.Name,
			Repositories:     repoName,
			State:            0,
			History:          history,
			UptimePercentage: -1,
			Frequency:        repoCfg.HealthCheck.Interval,
		}
	}

	latestHealthRow, err := s.repo.QueryLatestHealthByProject()
	if err != nil {
		return result, err
	}

	for repoName, row := range latestHealthRow {
		state := healthDTO.OverviewProjectStateUnmonitored
		if row.HttpStatus > 0 {
			if row.HttpStatus == http.StatusOK {
				result.AccuracyCount++
				state = healthDTO.OverviewProjectStateNormal
			} else {
				result.ErrorCount++
				state = healthDTO.OverviewProjectStateAbnormal
			}
		}
		// 处理每个项目信息
		project, ok := projects[repoName]
		if !ok {
			result.Projects = append(result.Projects, project)
			continue
		}
		project.State = state
		history, err := s.repo.QueryProjectHealthLast24Hours(repoName)
		if err == nil {
			project.History = history
			// 计算可用率
			sum := 0.0
			online := 0.0
			for _, v := range history {
				if v >= 0 {
					sum += v
					online += 100
				}
			}
			percentage := float64(sum) / float64(online) * 100
			project.UptimePercentage = math.Round(percentage*100) / 100
		}
		result.Projects = append(result.Projects, project)
	}
	// 部署记录
	webhookLogListPage, err := s.webhookSvc.QueryWebhookLogListPage(1, 10, nil)
	if err != nil {
		return result, err
	}
	for _, v := range webhookLogListPage.Data {
		result.Logs = append(result.Logs, healthDTO.OverviewLog{
			ID:       v.ID,
			Title:    v.Project,
			Execute:  v.Command,
			Datetime: v.CreatedAt,
		})
	}
	return result, nil
}

func (s Service) GetWebhookLogDetails(id string) (healthDTO.WebhookLogDetails, error) {
	var result healthDTO.WebhookLogDetails
	webhookLogDetails, err := s.webhookSvc.QueryWebhookLogDetails(id)
	if err != nil {
		return result, nil
	}
	result = healthDTO.WebhookLogDetails{
		ID:           webhookLogDetails.ID,
		Project:      webhookLogDetails.Project,
		Command:      webhookLogDetails.Command,
		Status:       webhookLogDetails.Status,
		StatusString: webhookLogDetails.StatusString,
		StartTime:    webhookLogDetails.StartTime,
		EndTime:      webhookLogDetails.EndTime,
		ExitCode:     webhookLogDetails.ExitCode,
		Error:        webhookLogDetails.Error,
		Stdout:       webhookLogDetails.Stdout,
		Stderr:       webhookLogDetails.Stderr,
	}
	return result, nil
}
