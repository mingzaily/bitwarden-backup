package scheduler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/bitwarden"
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// performBackupToDestinations 执行备份到多个目标
func (s *Scheduler) performBackupToDestinations(task database.BackupTask, backupLog *database.BackupLog) error {
	// 获取源服务器配置
	var sourceServer database.ServerConfig
	if err := database.DB.First(&sourceServer, task.SourceServerID).Error; err != nil {
		return fmt.Errorf("failed to get source server: %w", err)
	}

	// 创建 Bitwarden 客户端并导出数据
	client := bitwarden.NewClient()
	if err := client.ConfigServer(sourceServer.ServerURL); err != nil {
		return fmt.Errorf("failed to config server: %w", err)
	}

	if err := client.Login(sourceServer.ClientID, sourceServer.ClientSecret); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	if err := client.Unlock(sourceServer.MasterPassword); err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}

	// 生成临时备份文件
	timestamp := time.Now().Format("20060102_150405")
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("backup_%s_%s.json", task.Name, timestamp))

	if err := client.Export(tempFile, "json"); err != nil {
		client.Logout()
		return fmt.Errorf("failed to export: %w", err)
	}

	backupLog.BackupFile = tempFile

	// 备份到所有目标
	for _, dest := range task.Destinations {
		if !dest.Enabled {
			continue
		}

		if err := s.backupToDestination(dest, tempFile, task.Name, timestamp); err != nil {
			log.Printf("Failed to backup to destination %s: %v", dest.Name, err)
		}
	}

	client.Logout()
	return nil
}
