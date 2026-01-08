package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func (s *Scheduler) AddTask(task model.BackupTask) error {
	_, err := s.cron.AddFunc(task.CronExpression, func() {
		s.executeTask(task)
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}
	log.Printf("Task %s added with cron: %s", task.Name, task.CronExpression)
	return nil
}

func (s *Scheduler) executeTask(task model.BackupTask) {
	log.Printf("Executing task: %s", task.Name)

	startTime := time.Now()
	backupLog := model.BackupLog{
		TaskID:    task.ID,
		Status:    "running",
		StartTime: startTime,
	}
	database.DB.Create(&backupLog)

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
