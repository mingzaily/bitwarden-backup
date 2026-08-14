package model

import "time"

// LogEntry 单条执行日志
type LogEntry struct {
	Time    string `json:"time"`
	Message string `json:"message"`
}

// BackupLog 备份执行日志
type BackupLog struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	TaskID        uint       `gorm:"not null" json:"task_id"`
	Status        string     `gorm:"size:50;not null" json:"status"`
	Message       string     `gorm:"type:text" json:"message"`
	BackupFile    string     `gorm:"size:255" json:"backup_file"`
	ExecutionLogs string     `gorm:"type:text" json:"execution_logs"` // JSON 数组格式
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	CreatedAt     time.Time  `json:"created_at"`
}

// LogResponse is the safe API representation of a backup log. TaskName is
// loaded by the repository join so the UI can identify the execution without
// requiring one task query per log row.
type LogResponse struct {
	ID            uint       `json:"id"`
	TaskID        uint       `json:"task_id"`
	TaskName      string     `json:"task_name"`
	Status        string     `json:"status"`
	Message       string     `json:"message"`
	BackupFile    string     `json:"backup_file"`
	ExecutionLogs string     `json:"execution_logs"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	CreatedAt     time.Time  `json:"created_at"`
}
