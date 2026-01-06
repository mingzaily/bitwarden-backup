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

// getTempDir 获取临时目录，优先使用工作目录下的 .tmp
func getTempDir() string {
	// 尝试使用当前工作目录下的 .tmp 目录
	cwd, err := os.Getwd()
	if err == nil {
		tmpDir := filepath.Join(cwd, ".tmp")
		if err := os.MkdirAll(tmpDir, 0755); err == nil {
			return tmpDir
		}
	}
	// 回退到系统临时目录
	return os.TempDir()
}

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
	var tempFiles []string

	// 检查是否需要加密和非加密两种格式
	needEncrypted := false
	needPlain := false
	for _, dest := range task.Destinations {
		if !dest.Enabled {
			continue
		}
		if dest.Type == "local" || dest.Type == "webdav" {
			if dest.Encrypted {
				needEncrypted = true
			} else {
				needPlain = true
			}
		} else if dest.Type == "server" {
			needPlain = true
		}
	}

	// 导出非加密版本
	var plainFile string
	if needPlain {
		plainFile = filepath.Join(getTempDir(), fmt.Sprintf("backup_%s_%s.json", task.Name, timestamp))
		if err := client.Export(plainFile, "json"); err != nil {
			client.Logout()
			return fmt.Errorf("failed to export: %w", err)
		}
		tempFiles = append(tempFiles, plainFile)
	}

	// 导出加密版本
	var encryptedFile string
	if needEncrypted {
		encryptedFile = filepath.Join(getTempDir(), fmt.Sprintf("backup_%s_%s_encrypted.json", task.Name, timestamp))
		if err := client.Export(encryptedFile, "encrypted_json"); err != nil {
			client.Logout()
			return fmt.Errorf("failed to export encrypted: %w", err)
		}
		tempFiles = append(tempFiles, encryptedFile)
	}

	backupLog.BackupFile = plainFile
	if plainFile == "" {
		backupLog.BackupFile = encryptedFile
	}

	// 备份到所有目标
	for _, dest := range task.Destinations {
		if !dest.Enabled {
			continue
		}

		sourceFile := plainFile
		if (dest.Type == "local" || dest.Type == "webdav") && dest.Encrypted {
			sourceFile = encryptedFile
		}

		if err := s.backupToDestination(dest, sourceFile, task.Name, timestamp); err != nil {
			log.Printf("Failed to backup to destination %s: %v", dest.Name, err)
		}
	}

	// 清理临时文件
	for _, f := range tempFiles {
		if err := os.Remove(f); err != nil {
			log.Printf("Failed to remove temp file %s: %v", f, err)
		}
	}

	client.Logout()
	return nil
}
