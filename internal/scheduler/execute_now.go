package scheduler

import (
	"github.com/mingzaily/bitwarden-backup/internal/logger"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func (s *Scheduler) ExecuteTaskNow(task model.BackupTask) {
	logger.Module(logger.ModuleScheduler).Info("Manually executing task", "name", task.Name)

	startTime := time.Now()
	backupLog := model.BackupLog{
		TaskID:    task.ID,
		Status:    "running",
		StartTime: startTime,
	}
	database.DB.Create(&backupLog)

	if err := s.performBackup(task, &backupLog); err != nil {
		logger.Module(logger.ModuleScheduler).Error("Task failed", "name", task.Name, "error", err)
		endTime := time.Now()
		backupLog.Status = "failed"
		backupLog.Message = err.Error()
		backupLog.EndTime = &endTime
		database.DB.Save(&backupLog)
		return
	}

	endTime := time.Now()
	backupLog.Status = "success"
	backupLog.Message = "Backup completed successfully"
	backupLog.EndTime = &endTime
	database.DB.Save(&backupLog)
	logger.Module(logger.ModuleScheduler).Info("Task completed successfully", "name", task.Name)
}
