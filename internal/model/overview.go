package model

import "time"

// OverviewCount summarizes the total and enabled records for a resource.
type OverviewCount struct {
	Total   int64 `json:"total"`
	Enabled int64 `json:"enabled"`
}

// OverviewTaskStats summarizes the configured backup tasks.
type OverviewTaskStats struct {
	Total     int64 `json:"total"`
	Enabled   int64 `json:"enabled"`
	Scheduled int64 `json:"scheduled"`
}

// OverviewLogStats summarizes recent backup execution activity.
type OverviewLogStats struct {
	Total      int64 `json:"total"`
	Success24h int64 `json:"success_24h"`
	Failed24h  int64 `json:"failed_24h"`
	Running24h int64 `json:"running_24h"`
}

// OverviewTaskSummary is the non-sensitive task data used by the dashboard.
type OverviewTaskSummary struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Enabled          bool      `json:"enabled"`
	CronExpression   string    `json:"cron_expression"`
	SourceServerName string    `json:"source_server_name"`
	DestinationCount int       `json:"destination_count"`
	CreatedAt        time.Time `json:"created_at"`
}

// OverviewLogSummary is the non-sensitive log data used by the dashboard.
type OverviewLogSummary struct {
	ID         uint      `json:"id"`
	TaskID     uint      `json:"task_id"`
	TaskName   string    `json:"task_name"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
	BackupFile string    `json:"backup_file"`
	CreatedAt  time.Time `json:"created_at"`
}

// OverviewResponse is the dashboard payload. It intentionally contains only
// counts and display summaries, never decrypted connection credentials.
type OverviewResponse struct {
	Servers      OverviewCount         `json:"servers"`
	Destinations OverviewCount         `json:"destinations"`
	Tasks        OverviewTaskStats     `json:"tasks"`
	Logs         OverviewLogStats      `json:"logs"`
	RecentTasks  []OverviewTaskSummary `json:"recent_tasks"`
	RecentLogs   []OverviewLogSummary  `json:"recent_logs"`
}
