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

func (s *TaskService) Update(id uint, task *model.BackupTask) error {
	task.ID = id
	return s.repo.Update(task)
}

func (s *TaskService) Delete(id uint) error {
	return s.repo.Delete(id)
}
