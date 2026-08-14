package scheduler

import (
	"github.com/mingzaily/bitwarden-backup/internal/logger"

	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func (s *Scheduler) ExecuteTaskNow(task model.BackupTask) {
	// Manual execution must use the same bounded queue as scheduled execution.
	// This prevents callers from creating an unbounded goroutine per request.
	if task.ID == 0 || !s.TriggerTask(task.ID) {
		logger.Module(logger.ModuleScheduler).Warn("Manual task execution was not queued", "id", task.ID, "name", task.Name)
	}
}
