package scheduler

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/provider"
)

func (s *Scheduler) backupToDestination(dest model.BackupDestination, sourceFile, taskName, timestamp string) error {
	registry := provider.GetRegistry()

	p, err := registry.Get(dest.Type)
	if err != nil {
		return err
	}

	ctx := provider.BackupContext{
		SourceFile:  sourceFile,
		TaskName:    taskName,
		Timestamp:   timestamp,
		Destination: dest,
	}

	return p.Backup(ctx)
}
