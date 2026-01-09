package service

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetAll() ([]model.BackupTask, error) {
	return s.repo.FindAll()
}

func (s *TaskService) GetByID(id uint) (*model.BackupTask, error) {
	return s.repo.FindByID(id)
}

func (s *TaskService) GetEnabled() ([]model.BackupTask, error) {
	return s.repo.FindEnabled()
}

func (s *TaskService) Create(task *model.BackupTask) error {
	return s.repo.Create(task)
}

func (s *TaskService) CreateWithDestinations(task *model.BackupTask, destinationIDs []uint) error {
	return s.repo.CreateWithDestinations(task, destinationIDs)
}

func (s *TaskService) Update(id uint, task *model.BackupTask) error {
	task.ID = id
	return s.repo.Update(task)
}

func (s *TaskService) UpdateEnabled(id uint, enabled bool) error {
	return s.repo.UpdateEnabled(id, enabled)
}

func (s *TaskService) UpdateWithDestinations(id uint, task *model.BackupTask, destinationIDs []uint) error {
	task.ID = id
	return s.repo.UpdateWithDestinations(task, destinationIDs)
}

func (s *TaskService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetPaginated 分页获取任务
func (s *TaskService) GetPaginated(params model.PaginationParams) ([]model.BackupTask, int64, error) {
	return s.repo.FindPaginated(params)
}
