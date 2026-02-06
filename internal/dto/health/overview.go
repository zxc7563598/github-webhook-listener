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

type OverviewProject struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Repositories     string               `json:"repositories"`
	State            OverviewProjectState `json:"state"`
	History          []float64            `json:"history"`
	UptimePercentage float64              `json:"uptimePercentage"`
	Frequency        int                  `json:"frequency"`
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
