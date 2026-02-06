package health

import (
	"fmt"
	"math"
	"time"

	healthDTO "github.com/zxc7563598/github-webhook-listener/internal/dto/health"
	"github.com/zxc7563598/github-webhook-listener/internal/model"
)

type HealthRepository interface {
	HealthMonitoringCreate(project, errMessage string, httpStatus int, responseTimeMs int64) error
	QueryLatestHealthByProject() (map[string]healthDTO.LatestHealthRow, error)
	QueryProjectHealthLast24Hours(project string) ([]float64, error)
}

// 记录项目运行信息
func (r *gormRepo) HealthMonitoringCreate(project, errMessage string, httpStatus int, responseTimeMs int64) error {
	err := r.db.Create(&model.HealthMonitoring{
		Project:        project,
		Error:          errMessage,
		HttpStatus:     httpStatus,
		ResponseTimeMs: responseTimeMs,
	}).Error
	if err != nil {
		return fmt.Errorf("健康监控信息记录失败: %v", err)
	}
	return nil
}

// 获取每个项目最后一条监控信息
func (r *gormRepo) QueryLatestHealthByProject() (map[string]healthDTO.LatestHealthRow, error) {
	var rows []healthDTO.LatestHealthRow
	// 查询数据
	subQuery := r.db.Model(&model.HealthMonitoring{}).
		Select("project, MAX(created_at_unix) as max_created_at_unix").
		Group("project")

	err := r.db.Table("health_monitorings as h").
		Select("h.project, h.http_status, h.response_time_ms").
		Joins(`
			JOIN (?) as latest
			ON h.project = latest.project
			AND h.created_at_unix = latest.max_created_at_unix
		`, subQuery).Scan(&rows).Error

	if err != nil {
		return nil, fmt.Errorf("获取每个项目最后一条监控信息失败: %v", err)
	}
	// 整理数据
	result := make(map[string]healthDTO.LatestHealthRow, len(rows))
	for _, row := range rows {
		result[row.Project] = row
	}
	return result, nil
}

// 获取某个项目最近 24 小时的成功率（按小时）
func (r *gormRepo) QueryProjectHealthLast24Hours(project string) ([]float64, error) {
	now := time.Now().UTC()

	// 对齐到当前小时整点
	currentHour := now.Truncate(time.Hour)

	// 24 小时前
	startHour := currentHour.Add(-24 * time.Hour)

	result := make([]float64, 24)

	var stats []healthDTO.HourlyHealthStat
	err := r.db.Model(&model.HealthMonitoring{}).
		Select(`
			(created_at_unix / 3600) * 3600 AS hour_unix,
			COUNT(*) AS total_count,
			SUM(CASE WHEN http_status = 200 THEN 1 ELSE 0 END) AS success_count
		`).
		Where("project = ? AND created_at_unix >= ?", project, startHour.Unix()).
		Group("hour_unix").
		Order("hour_unix").
		Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("项目 %s 最近 24 小时数据查询失败: %v", project, err)
	}

	// hour_unix → index
	hourIndex := make(map[int64]int)
	for i := 0; i < 24; i++ {
		hourUnix := startHour.Add(time.Duration(i+1) * time.Hour).Unix()
		hourIndex[hourUnix] = i
	}

	// 填充结果
	for _, stat := range stats {
		idx, ok := hourIndex[stat.HourUnix]
		if !ok || stat.TotalCount == 0 {
			continue
		}

		percentage := float64(stat.SuccessCount) / float64(stat.TotalCount) * 100
		result[idx] = math.Round(percentage*100) / 100

		fmt.Printf("索引:%d, 最终结果: %.2f \n", idx, result[idx])
	}

	return result, nil
}
