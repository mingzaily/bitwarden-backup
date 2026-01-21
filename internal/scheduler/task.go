package scheduler

import (
	"fmt"

	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/logger"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func (s *Scheduler) LoadTasks() error {
	var tasks []model.BackupTask
	if err := database.DB.Preload("Destinations").Where("enabled = ?", true).Find(&tasks).Error; err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	scheduledCount := 0
	manualCount := 0

	for _, task := range tasks {
		if task.CronExpression == "" {
			logger.Module(logger.ModuleScheduler).Debug("Skipping manual task", "task", task.Name, "reason", "no cron expression")
			manualCount++
			continue
		}

		if err := s.AddTask(task); err != nil {
			logger.Module(logger.ModuleScheduler).Error("Failed to add task", "task", task.Name, "error", err)
		} else {
			scheduledCount++
		}
	}

	logger.Module(logger.ModuleScheduler).Info("Tasks loaded", "scheduled", scheduledCount, "manual", manualCount)
	return nil
}
