package scheduler

import (
	"fmt"
	"log"

	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// LoadTasks 从数据库加载所有启用的任务
func (s *Scheduler) LoadTasks() error {
	var tasks []database.BackupTask
	// 预加载关联的备份目标
	if err := database.DB.Preload("Destinations").Where("enabled = ?", true).Find(&tasks).Error; err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	scheduledCount := 0
	manualCount := 0

	for _, task := range tasks {
		// 跳过没有 Cron 表达式的任务（手动触发任务）
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
