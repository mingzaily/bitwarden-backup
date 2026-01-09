package scheduler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/bitwarden"
	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func getTempDir() string {
	cwd, err := os.Getwd()
	if err == nil {
		tmpDir := filepath.Join(cwd, ".tmp")
		if err := os.MkdirAll(tmpDir, 0755); err == nil {
			return tmpDir
		}
	}
	return os.TempDir()
}

func (s *Scheduler) performBackupToDestinations(task model.BackupTask, backupLog *model.BackupLog) error {
	var sourceServer model.ServerConfig
	if err := database.DB.First(&sourceServer, task.SourceServerID).Error; err != nil {
		return fmt.Errorf("failed to get source server: %w", err)
	}

	client := bitwarden.NewClient()

	// 使用 defer 确保无论成功还是失败都保存执行日志
	defer func() {
		if logs := client.GetLogs(); len(logs) > 0 {
			if logsJSON, err := json.Marshal(logs); err == nil {
				backupLog.ExecutionLogs = string(logsJSON)
			}
		}
	}()

	client.AddLog(fmt.Sprintf("Executing task: %s", task.Name))
	_ = client.Logout()

	if err := client.ConfigServer(sourceServer.ServerURL); err != nil {
		return fmt.Errorf("failed to config server: %w", err)
	}

	if err := client.Login(sourceServer.ClientID, sourceServer.ClientSecret); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	if err := client.Unlock(sourceServer.MasterPassword); err != nil {
		// 检测登录状态损坏，尝试重新登录
		if _, ok := err.(*bitwarden.ErrNotLoggedIn); ok {
			log.Printf("Login state corrupted, retrying login...")
			_ = client.Logout()
			if err := client.Login(sourceServer.ClientID, sourceServer.ClientSecret); err != nil {
				return fmt.Errorf("failed to re-login: %w", err)
			}
			if err := client.Unlock(sourceServer.MasterPassword); err != nil {
				return fmt.Errorf("failed to unlock after re-login: %w", err)
			}
		} else {
			return fmt.Errorf("failed to unlock: %w", err)
		}
	}

	timestamp := time.Now().Format("20060102_150405")
	var tempFiles []string

	needEncrypted := false
	needPlain := false
	var encryptionPassword string
	for _, dest := range task.Destinations {
		if !dest.Enabled {
			continue
		}
		if dest.Type == "local" || dest.Type == "webdav" || dest.Type == "s3" {
			if dest.Encrypted {
				needEncrypted = true
				if encryptionPassword == "" && dest.EncryptionPassword != "" {
					encryptionPassword = dest.EncryptionPassword
				}
			} else {
				needPlain = true
			}
		} else if dest.Type == "server" {
			needPlain = true
		}
	}

	if needEncrypted && encryptionPassword == "" {
		return fmt.Errorf("encryption password is required for encrypted backup destinations")
	}

	var plainFile string
	if needPlain {
		plainFile = filepath.Join(getTempDir(), fmt.Sprintf("backup_%s_%s.json", task.Name, timestamp))
		if err := client.Export(plainFile, "json"); err != nil {
			client.Logout()
			return fmt.Errorf("failed to export: %w", err)
		}
		tempFiles = append(tempFiles, plainFile)
	}

	var encryptedFile string
	if needEncrypted {
		encryptedFile = filepath.Join(getTempDir(), fmt.Sprintf("backup_%s_%s_encrypted.json", task.Name, timestamp))
		if err := client.Export(encryptedFile, "encrypted_json", encryptionPassword); err != nil {
			client.Logout()
			return fmt.Errorf("failed to export encrypted: %w", err)
		}
		tempFiles = append(tempFiles, encryptedFile)
	}

	var backupPaths []string
	for _, dest := range task.Destinations {
		if !dest.Enabled {
			continue
		}

		sourceFile := plainFile
		if (dest.Type == "local" || dest.Type == "webdav" || dest.Type == "s3") && dest.Encrypted {
			sourceFile = encryptedFile
		}

		targetPath, err := s.backupToDestination(dest, sourceFile, task.Name, timestamp)
		if err != nil {
			log.Printf("Failed to backup to destination %s: %v", dest.Name, err)
		} else if targetPath != "" {
			backupPaths = append(backupPaths, targetPath)
		}
	}

	// 存储第一个成功的备份路径
	if len(backupPaths) > 0 {
		backupLog.BackupFile = backupPaths[0]
	}

	for _, f := range tempFiles {
		if err := os.Remove(f); err != nil {
			log.Printf("Failed to remove temp file %s: %v", f, err)
		}
	}

	client.Logout()

	return nil
}
