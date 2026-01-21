package scheduler

import (
	"fmt"
	"github.com/mingzaily/bitwarden-backup/internal/logger"
	"strings"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// normalizeCron 将 5 位 cron 表达式转换为 6 位格式
// 5 位格式: 分 时 日 月 周 (标准 cron)
// 6 位格式: 秒 分 时 日 月 周 (robfig/cron)
func normalizeCron(expr string) string {
	fields := strings.Fields(expr)
	if len(fields) == 5 {
		return "0 " + expr // 在前面加 "0" 秒
	}
	return expr
}

func (s *Scheduler) AddTask(task model.BackupTask) error {
	cronExpr := normalizeCron(task.CronExpression)
	taskID := task.ID // 只捕获任务 ID，执行时重新查询最新数据
	entryID, err := s.cron.AddFunc(cronExpr, func() {
		s.enqueueTask(taskID)
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	// 保存任务ID到entry ID的映射
	s.mu.Lock()
	s.taskEntries[task.ID] = entryID
	s.mu.Unlock()

	logger.Module(logger.ModuleScheduler).Info("Task added", "name", task.Name, "id", task.ID, "cron", task.CronExpression)
	return nil
}

func (s *Scheduler) enqueueTask(taskID uint) {
	if s.stopped.Load() {
		return
	}

	s.queueMu.Lock()
	defer s.queueMu.Unlock()

	if s.stopped.Load() {
		return
	}
	if s.queuedTasks[taskID] {
		logger.Module(logger.ModuleScheduler).Info("Task already queued, skipping", "id", taskID)
		return
	}

	select {
	case s.taskQueue <- taskID:
		s.queuedTasks[taskID] = true
		logger.Module(logger.ModuleScheduler).Info("Task enqueued", "id", taskID)
	default:
		logger.Module(logger.ModuleScheduler).Warn("Task queue full, dropping task", "id", taskID)
	}
}

func (s *Scheduler) processTask(taskID uint) {
	defer s.removeFromQueue(taskID)
	defer func() {
		if r := recover(); r != nil {
			logger.Module(logger.ModuleScheduler).Error("Task execution panic recovered", "id", taskID, "panic", r)
		}
	}()

	var latestTask model.BackupTask
	if err := database.DB.Preload("Destinations").First(&latestTask, taskID).Error; err != nil {
		logger.Module(logger.ModuleScheduler).Error("Failed to fetch task for execution", "id", taskID, "error", err)
		return
	}
	if !latestTask.Enabled {
		logger.Module(logger.ModuleScheduler).Info("Task is disabled, skipping execution", "id", taskID, "name", latestTask.Name)
		return
	}
	s.executeTask(latestTask)
}

func (s *Scheduler) removeFromQueue(taskID uint) {
	s.queueMu.Lock()
	delete(s.queuedTasks, taskID)
	s.queueMu.Unlock()
}

// RemoveTask 从调度器中移除任务
func (s *Scheduler) RemoveTask(taskID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, exists := s.taskEntries[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.taskEntries, taskID)
		logger.Module(logger.ModuleScheduler).Info("Task removed from scheduler", "taskID", taskID)
	}
}

// UpdateTask 更新调度器中的任务（先移除再添加）
func (s *Scheduler) UpdateTask(task model.BackupTask) error {
	// 先移除旧任务
	s.RemoveTask(task.ID)

	// 如果任务禁用或没有cron表达式，不重新添加
	if !task.Enabled || task.CronExpression == "" {
		logger.Module(logger.ModuleScheduler).Info("Task not scheduled", "name", task.Name, "id", task.ID, "enabled", task.Enabled, "cron", task.CronExpression)
		return nil
	}

	// 添加新任务
	return s.AddTask(task)
}

func (s *Scheduler) executeTask(task model.BackupTask) {
	logger.Module(logger.ModuleScheduler).Info("Executing task", "name", task.Name)

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
