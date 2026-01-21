package scheduler

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func (s *Scheduler) performBackup(task model.BackupTask, backupLog *model.BackupLog) error {
	return s.performBackupToDestinations(task, backupLog)
}
