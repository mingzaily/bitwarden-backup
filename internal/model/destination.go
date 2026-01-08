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
	TargetServerID *uint `json:"target_server_id"`

	// 加密选项
	Encrypted          bool   `gorm:"default:false" json:"encrypted"`
	EncryptionPassword string `gorm:"size:500" json:"encryption_password"`

	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeSave GORM 钩子：保存前加密敏感字段
func (d *BackupDestination) BeforeSave(tx *gorm.DB) error {
	if d.WebDAVPassword != "" {
		encrypted, err := crypto.Encrypt(d.WebDAVPassword)
		if err != nil {
			return err
		}
		d.WebDAVPassword = encrypted
	}
	if d.S3AccessKey != "" {
		encrypted, err := crypto.Encrypt(d.S3AccessKey)
		if err != nil {
			return err
		}
		d.S3AccessKey = encrypted
	}
	if d.S3SecretKey != "" {
		encrypted, err := crypto.Encrypt(d.S3SecretKey)
		if err != nil {
			return err
		}
		d.S3SecretKey = encrypted
	}
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
