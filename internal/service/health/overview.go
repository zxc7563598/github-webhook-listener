package health

import (
	"math"
	"net/http"
	"time"

	"github.com/zxc7563598/github-webhook-listener/internal/config"
	healthDTO "github.com/zxc7563598/github-webhook-listener/internal/dto/health"
)

func (s Service) GetOverview(repos map[string]*config.RepoConfig) (healthDTO.Overview, error) {
	var result healthDTO.Overview
	var projects = make(map[string]healthDTO.OverviewProject)
	var overviewHistory = make([]healthDTO.OverviewHistory, 0, 24)
	// 获取 24 小时之前的时间
	now := time.Now().UTC()
	currentHour := now.Truncate(time.Hour)
	startHour := currentHour.Add(-24 * time.Hour)
	// 获取配置过的项目信息
	for repoName, repoCfg := range repos {
		result.ProjectCount++
		// 定义项目数据
		history := make([]float64, 24)
		for i := range history {
			history[i] = -1
		}
		projects[repoName] = healthDTO.OverviewProject{
			ID:           repoName,
			Name:         repoCfg.Name,
			Repositories: repoName,
			State:        0,
			History:      overviewHistory,
			Frequency:    repoCfg.HealthCheck.Interval,
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
		history, err := s.repo.QueryProjectHealthLast24Hours(repoName, startHour)
		overviewHistory = make([]healthDTO.OverviewHistory, 24)
		if err == nil {
			// hour_unix → index
			hourIndex := make(map[int64]int)
			for i := 0; i < 24; i++ {
				hourUnix := startHour.Add(time.Duration(i+1) * time.Hour).Unix()
				hourIndex[hourUnix] = i
			}
			// 填充结果
			for _, stat := range history {
				idx, ok := hourIndex[stat.HourUnix]
				if !ok || stat.TotalCount == 0 {
					continue
				}
				overviewHistory[idx] = healthDTO.OverviewHistory{
					Hour:            time.Unix(stat.HourUnix, 0).Format("2006-01-02 15"),
					TotalCount:      stat.TotalCount,
					SuccessCount:    stat.SuccessCount,
					AverageResponse: math.Round(float64(stat.SuccessResponseTimeMs) / float64(stat.SuccessCount)),
				}
			}
		}
		project.History = overviewHistory
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
		ID:        webhookLogDetails.ID,
		Project:   webhookLogDetails.Project,
		Command:   webhookLogDetails.Command,
		Status:    webhookLogDetails.Status,
		StartTime: webhookLogDetails.StartTime,
		EndTime:   webhookLogDetails.EndTime,
		ExitCode:  webhookLogDetails.ExitCode,
		Error:     webhookLogDetails.Error,
		Stdout:    webhookLogDetails.Stdout,
		Stderr:    webhookLogDetails.Stderr,
	}
	return result, nil
}
