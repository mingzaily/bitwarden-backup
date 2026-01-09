package service

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/repository"
)

type DestinationService struct {
	repo *repository.DestinationRepository
}

func NewDestinationService(repo *repository.DestinationRepository) *DestinationService {
	return &DestinationService{repo: repo}
}

func (s *DestinationService) GetAll() ([]model.BackupDestination, error) {
	return s.repo.FindAll()
}

func (s *DestinationService) GetByID(id uint) (*model.BackupDestination, error) {
	return s.repo.FindByID(id)
}

func (s *DestinationService) Create(dest *model.BackupDestination) error {
	return s.repo.Create(dest)
}

func (s *DestinationService) Update(id uint, dest *model.BackupDestination) error {
	dest.ID = id
	return s.repo.Update(dest)
}

func (s *DestinationService) UpdateEnabled(id uint, enabled bool) error {
	return s.repo.UpdateEnabled(id, enabled)
}

func (s *DestinationService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *DestinationService) Toggle(id uint) error {
	dest, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	dest.Enabled = !dest.Enabled
	return s.repo.Update(dest)
}

// GetPaginated 分页获取备份目标
func (s *DestinationService) GetPaginated(params model.PaginationParams) ([]model.BackupDestination, int64, error) {
	return s.repo.FindPaginated(params)
}
