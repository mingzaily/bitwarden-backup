package scheduler

import (
	"log"

	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/provider"
)

func (s *Scheduler) backupToDestination(dest model.BackupDestination, sourceFile, taskName, timestamp string) (string, error) {
	registry := provider.GetRegistry()

	p, err := registry.Get(dest.Type)
	if err != nil {
		return "", err
	}

	ctx := provider.BackupContext{
		SourceFile:  sourceFile,
		TaskName:    taskName,
		Timestamp:   timestamp,
		Destination: dest,
	}

	targetPath, err := p.Backup(ctx)
	if err != nil {
		return "", err
	}

	// 备份成功后执行清理
	if dest.MaxBackupCount > 0 {
		if rp, ok := p.(provider.RetentionProvider); ok {
			deleted, cleanupErr := rp.Cleanup(dest, dest.MaxBackupCount)
			if cleanupErr != nil {
				log.Printf("Warning: cleanup failed for %s: %v", dest.Name, cleanupErr)
			} else if deleted > 0 {
				log.Printf("Cleaned up %d old backup(s) from %s", deleted, dest.Name)
			}
		}
	}

	return targetPath, nil
}
