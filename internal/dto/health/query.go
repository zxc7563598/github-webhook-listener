package health

type HourlyHealthStat struct {
	HourUnix     int64 `gorm:"column:hour_unix"`
	TotalCount   int64 `gorm:"column:total_count"`
	SuccessCount int64 `gorm:"column:success_count"`
}

type LatestHealthRow struct {
	Project        string `gorm:"column:project"`
	HttpStatus     int    `gorm:"column:http_status"`
	ResponseTimeMs int64  `gorm:"column:response_time_ms"`
}
