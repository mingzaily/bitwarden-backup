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

// BeforeSave GORM 钩子：保存前加密敏感字段（防止双重加密）
func (s *ServerConfig) BeforeSave(tx *gorm.DB) error {
	if s.ClientID != "" && !crypto.IsEncrypted(s.ClientID) {
		encrypted, err := crypto.Encrypt(s.ClientID)
		if err != nil {
			return err
		}
		s.ClientID = encrypted
	}
	if s.ClientSecret != "" && !crypto.IsEncrypted(s.ClientSecret) {
		encrypted, err := crypto.Encrypt(s.ClientSecret)
		if err != nil {
			return err
		}
		s.ClientSecret = encrypted
	}
	if s.MasterPassword != "" && !crypto.IsEncrypted(s.MasterPassword) {
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

// ToResponse 转换为响应结构（隐藏敏感字段）
func (s *ServerConfig) ToResponse() ServerResponse {
	return ServerResponse{
		ID:         s.ID,
		Name:       s.Name,
		ServerURL:  s.ServerURL,
		ClientID:   s.ClientID,
		IsOfficial: s.IsOfficial,
		Enabled:    s.Enabled,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}

// ServerResponse 服务器响应结构（隐藏敏感字段）
type ServerResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	ServerURL  string    `json:"server_url"`
	ClientID   string    `json:"client_id"`
	IsOfficial bool      `json:"is_official"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ServerRequest 服务器请求 DTO
type ServerRequest struct {
	Name           string `json:"name"`
	ServerURL      string `json:"server_url"`
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	MasterPassword string `json:"master_password"`
	IsOfficial     bool   `json:"is_official"`
	Enabled        *bool  `json:"enabled"`
}
