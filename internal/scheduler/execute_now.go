package scheduler

import (
	"log"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// ExecuteTaskNow 立即执行任务
func (s *Scheduler) ExecuteTaskNow(task database.BackupTask) {
	log.Printf("Manually executing task: %s", task.Name)

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
