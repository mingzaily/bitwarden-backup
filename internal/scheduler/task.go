package scheduler

import (
	"fmt"
	"log"

	"github.com/mingzaily/bitwarden-backup/internal/database"
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
			log.Printf("Skipping manual task: %s (no cron expression)", task.Name)
			manualCount++
			continue
		}

		if err := s.AddTask(task); err != nil {
			log.Printf("Failed to add task %s: %v", task.Name, err)
		} else {
			scheduledCount++
		}
	}

	log.Printf("Loaded %d scheduled tasks, %d manual tasks", scheduledCount, manualCount)
	return nil
}
