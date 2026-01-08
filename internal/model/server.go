package model

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
	ClientID       string    `gorm:"size:500" json:"client_id"`
	ClientSecret   string    `gorm:"size:500" json:"client_secret"`
	MasterPassword string    `gorm:"size:500" json:"master_password"`
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
