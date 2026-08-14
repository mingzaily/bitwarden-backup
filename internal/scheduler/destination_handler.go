package scheduler

import (
	"context"
	"fmt"

	"github.com/mingzaily/bitwarden-backup/internal/logger"

	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/provider"
)

func (s *Scheduler) backupToDestination(requestCtx context.Context, dest model.BackupDestination, sourceFile, taskName, timestamp, filenameTemplate string, log func(source, message string)) (string, error) {
	registry := provider.GetRegistry()

	p, err := registry.Get(dest.Type)
	if err != nil {
		if log != nil {
			log(dest.Type, "获取备份服务商失败: "+err.Error())
		}
		return "", err
	}

	ctx := provider.BackupContext{
		Context:          requestCtx,
		SourceFile:       sourceFile,
		TaskName:         taskName,
		Timestamp:        timestamp,
		FilenameTemplate: filenameTemplate,
		Destination:      dest,
		Log:              log,
	}

	targetPath, err := p.Backup(ctx)
	if err != nil {
		ctx.AddLog(dest.Type, "服务商执行失败: "+err.Error())
		return "", err
	}

	// 备份成功后执行清理
	if dest.MaxBackupCount > 0 {
		if rp, ok := p.(provider.RetentionProvider); ok {
			deleted, cleanupErr := rp.Cleanup(ctx, dest.MaxBackupCount)
			if cleanupErr != nil {
				logger.Module(logger.ModuleScheduler).Warn("Cleanup failed", "destination", dest.Name, "error", cleanupErr)
				ctx.AddLog(dest.Type, "清理旧备份失败: "+cleanupErr.Error())
				// The backup artifact exists, but the destination is not in the
				// configured retention state. Do not let the task be recorded as
				// a complete success when cleanup could not be applied.
				return targetPath, fmt.Errorf("failed to clean up old backups: %w", cleanupErr)
			} else if deleted > 0 {
				logger.Module(logger.ModuleScheduler).Info("Cleaned up old backups", "count", deleted, "destination", dest.Name)
				ctx.AddLog(dest.Type, "已清理旧备份: "+fmt.Sprintf("%d 个", deleted))
			}
		}
	}

	return targetPath, nil
}
