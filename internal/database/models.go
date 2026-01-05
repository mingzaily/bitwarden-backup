package database

import (
	"time"
)

// ServerConfig 存储 Bitwarden 服务器配置
type ServerConfig struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	ServerURL      string    `gorm:"size:255;not null" json:"server_url"`
	ClientID       string    `gorm:"size:255" json:"client_id"`
	ClientSecret   string    `gorm:"size:255" json:"client_secret"`
	MasterPassword string    `gorm:"size:255" json:"master_password"`
	IsOfficial     bool      `gorm:"default:false" json:"is_official"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BackupTask 备份任务配置
type BackupTask struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:100;not null" json:"name"`
	SourceServerID   uint      `gorm:"not null" json:"source_server_id"`
	CronExpression   string    `gorm:"size:100;not null" json:"cron_expression"`
	Enabled          bool      `gorm:"default:true" json:"enabled"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// 关联的备份目标（多对多关系）
	Destinations     []BackupDestination `gorm:"many2many:task_destinations;" json:"destinations"`
}

// BackupDestination 备份目标配置
type BackupDestination struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	Type           string    `gorm:"size:20;not null" json:"type"` // local, webdav, server

	// 本地存储配置
	LocalPath      string    `gorm:"size:255" json:"local_path"`

	// WebDAV 配置
	WebDAVURL      string    `gorm:"size:255" json:"webdav_url"`
	WebDAVUsername string    `gorm:"size:100" json:"webdav_username"`
	WebDAVPassword string    `gorm:"size:255" json:"webdav_password"`
	WebDAVPath     string    `gorm:"size:255" json:"webdav_path"`

	// 目标服务器配置
	TargetServerID *uint     `json:"target_server_id"`

	Enabled        bool      `gorm:"default:true" json:"enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BackupLog 备份执行日志
type BackupLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TaskID      uint      `gorm:"not null" json:"task_id"`
	Status      string    `gorm:"size:50;not null" json:"status"` // success, failed, running
	Message     string    `gorm:"type:text" json:"message"`
	BackupFile  string    `gorm:"size:255" json:"backup_file"`
	StartTime   time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
}
