package scheduler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func (s *Scheduler) backupToLocal(dest model.BackupDestination, sourceFile, taskName, timestamp string) error {
	if err := os.MkdirAll(dest.LocalPath, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	targetFile := filepath.Join(dest.LocalPath, fmt.Sprintf("backup_%s_%s.json", taskName, timestamp))

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
