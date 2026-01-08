package model

import "time"

// BackupTask 备份任务配置
type BackupTask struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	SourceServerID uint      `gorm:"not null" json:"source_server_id"`
	CronExpression string    `gorm:"size:100" json:"cron_expression"`
	Enabled        bool      `gorm:"default:true" json:"enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联
	SourceServer ServerConfig        `json:"source_server"`
	Destinations []BackupDestination `gorm:"many2many:task_destinations;" json:"destinations"`
}
