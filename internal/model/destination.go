package model

import (
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/crypto"
	"gorm.io/gorm"
)

// BackupDestination 备份目标配置
type BackupDestination struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;not null" json:"name"`
	Type string `gorm:"size:20;not null" json:"type"`

	// 本地存储配置
	LocalPath string `gorm:"size:255" json:"local_path"`

	// WebDAV 配置
	WebDAVURL      string `gorm:"size:255" json:"webdav_url"`
	WebDAVUsername string `gorm:"size:100" json:"webdav_username"`
	WebDAVPassword string `gorm:"size:500" json:"webdav_password"`
	WebDAVPath     string `gorm:"size:255" json:"webdav_path"`

	// S3 配置
	S3Endpoint  string `gorm:"size:255" json:"s3_endpoint"`
	S3Region    string `gorm:"size:100" json:"s3_region"`
	S3Bucket    string `gorm:"size:100" json:"s3_bucket"`
	S3AccessKey string `gorm:"size:500" json:"s3_access_key"`
	S3SecretKey string `gorm:"size:500" json:"s3_secret_key"`
	S3Path      string `gorm:"size:255" json:"s3_path"`

	// 目标服务器配置
	TargetServerID *uint         `json:"target_server_id"`
	TargetServer   *ServerConfig `gorm:"foreignKey:TargetServerID" json:"target_server,omitempty"`

	// 加密选项
	Encrypted          bool   `gorm:"default:false" json:"encrypted"`
	EncryptionPassword string `gorm:"size:500" json:"encryption_password"`

	// 备份保留策略
	MaxBackupCount int `gorm:"default:0" json:"max_backup_count"` // 0 表示不限制

	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeSave GORM 钩子：保存前加密敏感字段（防止双重加密）
func (d *BackupDestination) BeforeSave(tx *gorm.DB) error {
	if d.WebDAVPassword != "" && !crypto.IsEncrypted(d.WebDAVPassword) {
		encrypted, err := crypto.Encrypt(d.WebDAVPassword)
		if err != nil {
			return err
		}
		d.WebDAVPassword = encrypted
	}
	if d.S3AccessKey != "" && !crypto.IsEncrypted(d.S3AccessKey) {
		encrypted, err := crypto.Encrypt(d.S3AccessKey)
		if err != nil {
			return err
		}
		d.S3AccessKey = encrypted
	}
	if d.S3SecretKey != "" && !crypto.IsEncrypted(d.S3SecretKey) {
		encrypted, err := crypto.Encrypt(d.S3SecretKey)
		if err != nil {
			return err
		}
		d.S3SecretKey = encrypted
	}
	if d.EncryptionPassword != "" && !crypto.IsEncrypted(d.EncryptionPassword) {
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
	if d.WebDAVPassword != "" {
		decrypted, err := crypto.Decrypt(d.WebDAVPassword)
		if err != nil {
			return err
		}
		d.WebDAVPassword = decrypted
	}
	if d.S3AccessKey != "" {
		decrypted, err := crypto.Decrypt(d.S3AccessKey)
		if err != nil {
			return err
		}
		d.S3AccessKey = decrypted
	}
	if d.S3SecretKey != "" {
		decrypted, err := crypto.Decrypt(d.S3SecretKey)
		if err != nil {
			return err
		}
		d.S3SecretKey = decrypted
	}
	if d.EncryptionPassword != "" {
		decrypted, err := crypto.Decrypt(d.EncryptionPassword)
		if err != nil {
			return err
		}
		d.EncryptionPassword = decrypted
	}
	return nil
}

// DestinationResponse 备份目标响应 DTO（隐藏敏感数据）
type DestinationResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	LocalPath      string    `json:"local_path,omitempty"`
	WebDAVURL      string    `json:"webdav_url,omitempty"`
	WebDAVUsername string    `json:"webdav_username,omitempty"`
	WebDAVPath     string    `json:"webdav_path,omitempty"`
	S3Endpoint     string    `json:"s3_endpoint,omitempty"`
	S3Region       string    `json:"s3_region,omitempty"`
	S3Bucket       string    `json:"s3_bucket,omitempty"`
	S3AccessKey    string    `json:"s3_access_key,omitempty"`
	S3Path         string    `json:"s3_path,omitempty"`
	TargetServerID *uint     `json:"target_server_id,omitempty"`
	Encrypted      bool      `json:"encrypted"`
	MaxBackupCount int       `json:"max_backup_count"`
	Enabled        bool      `json:"enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DisplayPath    string    `json:"display_path"`
	TypeLabel      string    `json:"type_label"`
}

// maskSensitiveField 掩码敏感字段，只显示前4位和后4位
func maskSensitiveField(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

// ToResponse 转换为响应结构（隐藏敏感字段）
func (d *BackupDestination) ToResponse() DestinationResponse {
	// 对 S3AccessKey 进行掩码处理
	maskedS3AccessKey := ""
	if d.S3AccessKey != "" {
		maskedS3AccessKey = maskSensitiveField(d.S3AccessKey)
	}

	return DestinationResponse{
		ID:             d.ID,
		Name:           d.Name,
		Type:           d.Type,
		LocalPath:      d.LocalPath,
		WebDAVURL:      d.WebDAVURL,
		WebDAVUsername: d.WebDAVUsername,
		WebDAVPath:     d.WebDAVPath,
		S3Endpoint:     d.S3Endpoint,
		S3Region:       d.S3Region,
		S3Bucket:       d.S3Bucket,
		S3AccessKey:    maskedS3AccessKey,
		S3Path:         d.S3Path,
		TargetServerID: d.TargetServerID,
		Encrypted:      d.Encrypted,
		MaxBackupCount: d.MaxBackupCount,
		Enabled:        d.Enabled,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
		DisplayPath:    d.GetDisplayPath(),
		TypeLabel:      d.GetTypeLabel(),
	}
}

// GetDisplayPath 获取展示路径
func (d *BackupDestination) GetDisplayPath() string {
	switch d.Type {
	case "local":
		return d.LocalPath
	case "webdav":
		return d.WebDAVURL + d.WebDAVPath
	case "s3":
		return "s3://" + d.S3Bucket + d.S3Path
	case "server":
		if d.TargetServer != nil {
			return d.TargetServer.Name + " · " + d.TargetServer.ServerURL
		}
		return "目标服务器"
	default:
		return ""
	}
}

// GetTypeLabel 获取类型标签
func (d *BackupDestination) GetTypeLabel() string {
	labels := map[string]string{
		"local":  "本地存储",
		"webdav": "WebDAV",
		"s3":     "S3",
		"server": "服务器",
	}
	if label, ok := labels[d.Type]; ok {
		return label
	}
	return d.Type
}
