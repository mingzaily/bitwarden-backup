package handler

import (
	"github.com/mingzaily/bitwarden-backup/internal/repository"
	"github.com/mingzaily/bitwarden-backup/internal/service"
	"gorm.io/gorm"
)

var (
	serverSvc      *service.ServerService
	destinationSvc *service.DestinationService
	taskSvc        *service.TaskService
	logSvc         *service.LogService
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
