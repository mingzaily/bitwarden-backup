package scheduler

import (
	"fmt"

	"github.com/mingzaily/bitwarden-backup/internal/bitwarden"
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// performMigration 执行数据迁移到目标服务器
func (s *Scheduler) performMigration(client *bitwarden.Client, targetServerID uint, backupFile string) error {
	// 获取目标服务器配置
	var targetServer database.ServerConfig
	if err := database.DB.First(&targetServer, targetServerID).Error; err != nil {
		return fmt.Errorf("failed to get target server: %w", err)
	}

	// 登出当前服务器
	if err := client.Logout(); err != nil {
		return fmt.Errorf("failed to logout from source: %w", err)
	}

	// 配置目标服务器
	if err := client.ConfigServer(targetServer.ServerURL); err != nil {
		return fmt.Errorf("failed to config target server: %w", err)
	}

	// 登录目标服务器
	if err := client.Login(targetServer.ClientID, targetServer.ClientSecret); err != nil {
		return fmt.Errorf("failed to login to target: %w", err)
	}

	// 解锁目标服务器
	if err := client.Unlock(targetServer.MasterPassword); err != nil {
		return fmt.Errorf("failed to unlock target: %w", err)
	}

	// 导入数据
	if err := client.Import(backupFile, "json"); err != nil {
		return fmt.Errorf("failed to import: %w", err)
	}

	return nil
}
