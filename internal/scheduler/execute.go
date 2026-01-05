package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// AddTask 添加任务到调度器
func (s *Scheduler) AddTask(task database.BackupTask) error {
	_, err := s.cron.AddFunc(task.CronExpression, func() {
		s.executeTask(task)
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}
	log.Printf("Task %s added with cron: %s", task.Name, task.CronExpression)
	return nil
}

// executeTask 执行备份任务
func (s *Scheduler) executeTask(task database.BackupTask) {
	log.Printf("Executing task: %s", task.Name)

	startTime := time.Now()
	backupLog := database.BackupLog{
		TaskID:    task.ID,
		Status:    "running",
		StartTime: startTime,
	}
	database.DB.Create(&backupLog)

	// 执行备份
	if err := s.performBackup(task, &backupLog); err != nil {
		log.Printf("Task %s failed: %v", task.Name, err)
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
	log.Printf("Task %s completed successfully", task.Name)
}
