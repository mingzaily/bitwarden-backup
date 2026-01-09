package provider

import (
	"fmt"

	"github.com/mingzaily/bitwarden-backup/internal/bitwarden"
	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// ServerProvider 服务器备份提供者
type ServerProvider struct{}

// NewServerProvider 创建服务器备份提供者
func NewServerProvider() *ServerProvider {
	return &ServerProvider{}
}

// Type 返回提供者类型
func (p *ServerProvider) Type() string {
	return "server"
}

// Backup 执行服务器备份（导入到目标服务器），返回目标服务器信息
func (p *ServerProvider) Backup(ctx BackupContext) (string, error) {
	dest := ctx.Destination

	if dest.TargetServerID == nil {
		return "", fmt.Errorf("target server id is nil")
	}

	var targetServer model.ServerConfig
	if err := database.DB.First(&targetServer, *dest.TargetServerID).Error; err != nil {
		return "", fmt.Errorf("failed to get target server: %w", err)
	}

	client := bitwarden.NewClient()
	if err := client.ConfigServer(targetServer.ServerURL); err != nil {
		return "", fmt.Errorf("failed to config target server: %w", err)
	}

	if err := client.Login(targetServer.ClientID, targetServer.ClientSecret); err != nil {
		return "", fmt.Errorf("failed to login to target: %w", err)
	}

	if err := client.Unlock(targetServer.MasterPassword); err != nil {
		return "", fmt.Errorf("failed to unlock target: %w", err)
	}

	if err := client.Import(ctx.SourceFile, "json"); err != nil {
		client.Logout()
		return "", fmt.Errorf("failed to import: %w", err)
	}

	client.Logout()
	// 返回目标服务器信息
	return fmt.Sprintf("server://%s", targetServer.Name), nil
}
