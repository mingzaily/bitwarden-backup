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

	for _, task := range tasks {
		if err := s.AddTask(task); err != nil {
			log.Printf("Failed to add task %s: %v", task.Name, err)
		}
	}

	log.Printf("Loaded %d tasks", len(tasks))
	return nil
}
