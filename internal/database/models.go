package database

import (
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/crypto"
	"gorm.io/gorm"
)

// ServerConfig 存储 Bitwarden 服务器配置
type ServerConfig struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	ServerURL      string    `gorm:"size:255;not null" json:"server_url"`
	ClientID       string    `gorm:"size:500" json:"client_id"`       // 增加长度以容纳加密数据
	ClientSecret   string    `gorm:"size:500" json:"client_secret"`   // 增加长度以容纳加密数据
	MasterPassword string    `gorm:"size:500" json:"master_password"` // 增加长度以容纳加密数据
	IsOfficial     bool      `gorm:"default:false" json:"is_official"`
	Enabled        bool      `gorm:"default:true" json:"enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BeforeSave GORM 钩子：保存前加密敏感字段
func (s *ServerConfig) BeforeSave(tx *gorm.DB) error {
	if s.ClientID != "" {
		encrypted, err := crypto.Encrypt(s.ClientID)
		if err != nil {
			return err
		}
		s.ClientID = encrypted
	}

	if s.ClientSecret != "" {
		encrypted, err := crypto.Encrypt(s.ClientSecret)
		if err != nil {
			return err
		}
		s.ClientSecret = encrypted
	}

	if s.MasterPassword != "" {
		encrypted, err := crypto.Encrypt(s.MasterPassword)
		if err != nil {
			return err
		}
		s.MasterPassword = encrypted
	}

	return nil
}

// AfterFind GORM 钩子：查询后解密敏感字段
func (s *ServerConfig) AfterFind(tx *gorm.DB) error {
	if s.ClientID != "" {
		decrypted, err := crypto.Decrypt(s.ClientID)
		if err != nil {
			return err
		}
		s.ClientID = decrypted
	}

	if s.ClientSecret != "" {
		decrypted, err := crypto.Decrypt(s.ClientSecret)
		if err != nil {
			return err
		}
		s.ClientSecret = decrypted
	}

	if s.MasterPassword != "" {
		decrypted, err := crypto.Decrypt(s.MasterPassword)
		if err != nil {
			return err
		}
		s.MasterPassword = decrypted
	}

	return nil
}

// BackupTask 备份任务配置
type BackupTask struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:100;not null" json:"name"`
	SourceServerID   uint      `gorm:"not null" json:"source_server_id"`
	CronExpression   string    `gorm:"size:100" json:"cron_expression"` // 可选，为空则仅支持手动触发
	Enabled          bool      `gorm:"default:true" json:"enabled"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// 关联的源服务器
	SourceServer     ServerConfig `json:"source_server"`

	// 关联的备份目标（多对多关系）
	Destinations     []BackupDestination `gorm:"many2many:task_destinations;" json:"destinations"`
}

// BackupDestination 备份目标配置
type BackupDestination struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	Type           string    `gorm:"size:20;not null" json:"type"` // local, webdav, server, s3

	// 本地存储配置
	LocalPath      string    `gorm:"size:255" json:"local_path"`

	// WebDAV 配置
	WebDAVURL      string    `gorm:"size:255" json:"webdav_url"`
	WebDAVUsername string    `gorm:"size:100" json:"webdav_username"`
	WebDAVPassword string    `gorm:"size:500" json:"webdav_password"` // 增加长度以容纳加密数据
	WebDAVPath     string    `gorm:"size:255" json:"webdav_path"`

	// S3 配置
	S3Endpoint     string    `gorm:"size:255" json:"s3_endpoint"`
	S3Region       string    `gorm:"size:100" json:"s3_region"`
	S3Bucket       string    `gorm:"size:100" json:"s3_bucket"`
	S3AccessKey    string    `gorm:"size:500" json:"s3_access_key"`    // 增加长度以容纳加密数据
	S3SecretKey    string    `gorm:"size:500" json:"s3_secret_key"`    // 增加长度以容纳加密数据
	S3Path         string    `gorm:"size:255" json:"s3_path"`

	// 目标服务器配置
	TargetServerID *uint     `json:"target_server_id"`

	// 加密选项（仅本地、WebDAV 和 S3 有效）
	Encrypted           bool   `gorm:"default:false" json:"encrypted"`
	EncryptionPassword  string `gorm:"size:500" json:"encryption_password"` // 加密备份的密码

	Enabled        bool      `gorm:"default:true" json:"enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BeforeSave GORM 钩子：保存前加密敏感字段
func (d *BackupDestination) BeforeSave(tx *gorm.DB) error {
	// 加密 WebDAV 密码
	if d.WebDAVPassword != "" {
		encrypted, err := crypto.Encrypt(d.WebDAVPassword)
		if err != nil {
			return err
		}
		d.WebDAVPassword = encrypted
	}

	// 加密 S3 AccessKey
	if d.S3AccessKey != "" {
		encrypted, err := crypto.Encrypt(d.S3AccessKey)
		if err != nil {
			return err
		}
		d.S3AccessKey = encrypted
	}

	// 加密 S3 SecretKey
	if d.S3SecretKey != "" {
		encrypted, err := crypto.Encrypt(d.S3SecretKey)
		if err != nil {
			return err
		}
		d.S3SecretKey = encrypted
	}

	// 加密 EncryptionPassword
	if d.EncryptionPassword != "" {
		encrypted, err := crypto.Encrypt(d.EncryptionPassword)
		if err != nil {
			return err
		}
		d.EncryptionPassword = encrypted
	}

	return nil
}

// AfterFind GORM 钩子：查询后解密敏感字段
func (d *BackupDestination) AfterFind(tx *gorm.DB) error {
	// 解密 WebDAV 密码
	if d.WebDAVPassword != "" {
		decrypted, err := crypto.Decrypt(d.WebDAVPassword)
		if err != nil {
			return err
		}
		d.WebDAVPassword = decrypted
	}

	// 解密 S3 AccessKey
	if d.S3AccessKey != "" {
		decrypted, err := crypto.Decrypt(d.S3AccessKey)
		if err != nil {
			return err
		}
		d.S3AccessKey = decrypted
	}

	// 解密 S3 SecretKey
	if d.S3SecretKey != "" {
		decrypted, err := crypto.Decrypt(d.S3SecretKey)
		if err != nil {
			return err
		}
		d.S3SecretKey = decrypted
	}

	// 解密 EncryptionPassword
	if d.EncryptionPassword != "" {
		decrypted, err := crypto.Decrypt(d.EncryptionPassword)
		if err != nil {
			return err
		}
		d.EncryptionPassword = decrypted
	}

	return nil
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
