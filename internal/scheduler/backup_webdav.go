package scheduler

import (
	"fmt"
	"path/filepath"

	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/webdav"
)

// backupToWebDAV 备份到 WebDAV
func (s *Scheduler) backupToWebDAV(dest database.BackupDestination, sourceFile, taskName, timestamp string) error {
	// 创建 WebDAV 客户端
	client := webdav.NewClient(dest.WebDAVURL, dest.WebDAVUsername, dest.WebDAVPassword)

	// 生成远程文件路径
	remoteFile := filepath.Join(dest.WebDAVPath, fmt.Sprintf("backup_%s_%s.json", taskName, timestamp))

	// 上传文件
	if err := client.UploadFile(sourceFile, remoteFile); err != nil {
		return fmt.Errorf("failed to upload to webdav: %w", err)
	}

	return nil
}
