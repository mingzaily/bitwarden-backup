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
	TriggerTask(taskID uint) bool
}

// ServerService describes the server operations needed by the HTTP layer.
// Keeping the dependency as an interface makes handlers independently
// testable and avoids package-level mutable state.
type ServerService interface {
	GetByID(id uint) (*model.ServerConfig, error)
	GetPaginated(params model.PaginationParams, enabled *bool) ([]model.ServerConfig, int64, error)
	Create(server *model.ServerConfig) error
	Update(id uint, server *model.ServerConfig) error
	UpdateEnabled(id uint, enabled bool) error
	Delete(id uint) error
}

// DestinationService describes the destination operations needed by handlers.
type DestinationService interface {
	GetByID(id uint) (*model.BackupDestination, error)
	GetPaginated(params model.PaginationParams) ([]model.BackupDestination, int64, error)
	Create(destination *model.BackupDestination) error
	Update(id uint, destination *model.BackupDestination) error
	UpdateEnabled(id uint, enabled bool) error
	Delete(id uint) error
	Toggle(id uint) error
	TestConnection(id uint) error
}

// TaskService describes the task operations needed by handlers.
type TaskService interface {
	GetByID(id uint) (*model.BackupTask, error)
	GetPaginated(params model.PaginationParams) ([]model.BackupTask, int64, error)
	CreateWithDestinations(task *model.BackupTask, destinationIDs []uint) error
	UpdateEnabled(id uint, enabled bool) error
	UpdateWithDestinations(id uint, task *model.BackupTask, destinationIDs []uint) error
	Delete(id uint) error
}

// LogService describes the log operations needed by handlers.
type LogService interface {
	GetPaginated(params model.PaginationParams, taskID *uint) ([]model.LogResponse, int64, error)
	DeleteByIDs(ids []uint) (int64, error)
}

// OverviewService describes the aggregate data needed by the dashboard.
type OverviewService interface {
	Get() (model.OverviewResponse, error)
}

// API owns all HTTP-layer dependencies. One instance is created during
// startup and passed to the router; handlers no longer rely on global state.
type API struct {
	serverService      ServerService
	destinationService DestinationService
	taskService        TaskService
	logService         LogService
	overviewService    OverviewService
	scheduler          TaskScheduler
}

// New constructs the HTTP API with the application's default services.
func New(db *gorm.DB) *API {
	api := NewWithDependencies(
		service.NewServerService(repository.NewServerRepository(db)),
		service.NewDestinationService(repository.NewDestinationRepository(db)),
		service.NewTaskService(repository.NewTaskRepository(db)),
		service.NewLogService(repository.NewLogRepository(db)),
		nil,
	)
	api.SetOverviewService(service.NewOverviewService(repository.NewOverviewRepository(db)))
	return api
}

// NewWithDependencies constructs an API with explicit dependencies. It is
// useful for focused handler tests and for embedding the application.
func NewWithDependencies(
	serverService ServerService,
	destinationService DestinationService,
	taskService TaskService,
	logService LogService,
	scheduler TaskScheduler,
) *API {
	return &API{
		serverService:      serverService,
		destinationService: destinationService,
		taskService:        taskService,
		logService:         logService,
		scheduler:          scheduler,
	}
}

// SetScheduler sets the scheduler used by task mutation and execution
// handlers. The setter is kept small because the scheduler starts after the
// API is constructed during application startup.
func (a *API) SetScheduler(scheduler TaskScheduler) {
	a.scheduler = scheduler
}

// SetOverviewService injects the aggregate read service used by the home
// dashboard. It remains a setter so existing handler tests and embedders that
// use NewWithDependencies do not need to construct an overview dependency.
func (a *API) SetOverviewService(overviewService OverviewService) {
	a.overviewService = overviewService
}
