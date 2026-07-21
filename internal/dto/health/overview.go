package health

import "github.com/zxc7563598/github-webhook-listener/internal/model"

type OverviewProjectState int

const (
	OverviewProjectStateUnmonitored OverviewProjectState = iota
	OverviewProjectStateNormal
	OverviewProjectStateAbnormal
)

func (state OverviewProjectState) String() string {
	switch state {
	case OverviewProjectStateUnmonitored:
		return "未监听"
	case OverviewProjectStateNormal:
		return "运行中"
	case OverviewProjectStateAbnormal:
		return "异常"
	default:
		return "未知"
	}
}

type OverviewHistory struct {
	Hour            string  `json:"hour"`
	TotalCount      int64   `json:"total_count"`
	SuccessCount    int64   `json:"success_count"`
	AverageResponse float64 `json:"average_response"`
}

type OverviewProject struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Repositories string               `json:"repositories"`
	State        OverviewProjectState `json:"state"`
	History      []OverviewHistory    `json:"history"`
	Frequency    int                  `json:"frequency"`
}

type OverviewLog struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Execute  string `json:"execute"`
	Datetime string `json:"datetime"`
}

type Overview struct {
	ProjectCount  int               `json:"projectCount"`
	AccuracyCount int               `json:"accuracyCount"`
	ErrorCount    int               `json:"errorCount"`
	Projects      []OverviewProject `json:"projects"`
	Logs          []OverviewLog     `json:"logs"`
}

type WebhookLogDetails struct {
	ID        string                 `json:"id"`
	Project   string                 `json:"project"`
	Command   string                 `json:"command"`
	Status    model.WebhookLogStatus `json:"status"`
	StartTime string                 `json:"startTime"`
	EndTime   string                 `json:"endTime"`
	ExitCode  int                    `json:"exitCode"`
	Error     string                 `json:"error"`
	Stdout    string                 `json:"stdout"`
	Stderr    string                 `json:"stderr"`
}

type GetWebhookLogDetailsRequest struct {
	ID string `json:"id" form:"id"`
}
