package scheduler

import (
	"fmt"

	"github.com/mingzaily/bitwarden-backup/internal/bitwarden"
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// backupToServer 备份到目标服务器（导入）
func (s *Scheduler) backupToServer(dest database.BackupDestination, sourceFile string) error {
	if dest.TargetServerID == nil {
		return fmt.Errorf("target server id is nil")
	}

	// 获取目标服务器配置
	var targetServer database.ServerConfig
	if err := database.DB.First(&targetServer, *dest.TargetServerID).Error; err != nil {
		return fmt.Errorf("failed to get target server: %w", err)
	}

	// 创建新的客户端连接到目标服务器
	client := bitwarden.NewClient()
	if err := client.ConfigServer(targetServer.ServerURL); err != nil {
		return fmt.Errorf("failed to config target server: %w", err)
	}

	if err := client.Login(targetServer.ClientID, targetServer.ClientSecret); err != nil {
		return fmt.Errorf("failed to login to target: %w", err)
	}

	if err := client.Unlock(targetServer.MasterPassword); err != nil {
		return fmt.Errorf("failed to unlock target: %w", err)
	}

	// 导入数据
	if err := client.Import(sourceFile, "json"); err != nil {
		client.Logout()
		return fmt.Errorf("failed to import: %w", err)
	}

	client.Logout()
	return nil
}
