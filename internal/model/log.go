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
