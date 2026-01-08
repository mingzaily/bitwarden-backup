package scheduler

import (
	"fmt"
	"path/filepath"

	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/webdav"
)

func (s *Scheduler) backupToWebDAV(dest model.BackupDestination, sourceFile, taskName, timestamp string) error {
	client := webdav.NewClient(dest.WebDAVURL, dest.WebDAVUsername, dest.WebDAVPassword)
	remoteFile := filepath.Join(dest.WebDAVPath, fmt.Sprintf("backup_%s_%s.json", taskName, timestamp))

	if err := client.UploadFile(sourceFile, remoteFile); err != nil {
		return fmt.Errorf("failed to upload to webdav: %w", err)
	}

	return nil
}
