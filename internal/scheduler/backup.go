package scheduler

import (
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// performBackup 执行实际的备份操作
func (s *Scheduler) performBackup(task database.BackupTask, backupLog *database.BackupLog) error {
	// 使用新的多目标备份逻辑
	return s.performBackupToDestinations(task, backupLog)
}
