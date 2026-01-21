package handler

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/repository"
	"github.com/mingzaily/bitwarden-backup/internal/service"
	"gorm.io/gorm"
)

// TaskScheduler 任务调度器接口，用于动态更新调度任务
type TaskScheduler interface {
	AddTask(task model.BackupTask) error
	RemoveTask(taskID uint)
	UpdateTask(task model.BackupTask) error
}

var (
	serverSvc      *service.ServerService
	destinationSvc *service.DestinationService
	taskSvc        *service.TaskService
	logSvc         *service.LogService
	taskScheduler  TaskScheduler // 任务调度器实例
)

// Init 初始化所有 handler 依赖
func Init(db *gorm.DB) {
	// 初始化 Repository
	serverRepo := repository.NewServerRepository(db)
	destRepo := repository.NewDestinationRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	logRepo := repository.NewLogRepository(db)

	// 初始化 Service
	serverSvc = service.NewServerService(serverRepo)
	destinationSvc = service.NewDestinationService(destRepo)
	taskSvc = service.NewTaskService(taskRepo)
	logSvc = service.NewLogService(logRepo)
}

// SetScheduler 设置任务调度器实例
func SetScheduler(s TaskScheduler) {
	taskScheduler = s
}
