package scheduler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// backupToLocal 备份到本地存储
func (s *Scheduler) backupToLocal(dest database.BackupDestination, sourceFile, taskName, timestamp string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(dest.LocalPath, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	// 生成目标文件名
	targetFile := filepath.Join(dest.LocalPath, fmt.Sprintf("backup_%s_%s.json", taskName, timestamp))

	// 复制文件
	source, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	target, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
