package service

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/repository"
)

type LogService struct {
	repo *repository.LogRepository
}

func NewLogService(repo *repository.LogRepository) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) GetAll() ([]model.BackupLog, error) {
	return s.repo.FindAll()
}

func (s *LogService) GetByTaskID(taskID uint) ([]model.BackupLog, error) {
	return s.repo.FindByTaskID(taskID)
}

// GetPaginated 分页获取日志
func (s *LogService) GetPaginated(params model.PaginationParams, taskID *uint) ([]model.BackupLog, int64, error) {
	return s.repo.FindPaginated(params, taskID)
}
