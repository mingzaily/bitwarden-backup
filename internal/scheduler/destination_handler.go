package scheduler

import (
	"context"

	"github.com/mingzaily/bitwarden-backup/internal/logger"

	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/provider"
)

func (s *Scheduler) backupToDestination(requestCtx context.Context, dest model.BackupDestination, sourceFile, taskName, timestamp string) (string, error) {
	registry := provider.GetRegistry()

	p, err := registry.Get(dest.Type)
	if err != nil {
		return "", err
	}

	ctx := provider.BackupContext{
		Context:     requestCtx,
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
				logger.Module(logger.ModuleScheduler).Warn("Cleanup failed", "destination", dest.Name, "error", cleanupErr)
			} else if deleted > 0 {
				logger.Module(logger.ModuleScheduler).Info("Cleaned up old backups", "count", deleted, "destination", dest.Name)
			}
		}
	}

	return targetPath, nil
}
