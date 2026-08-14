package service

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/repository"
)

type OverviewService struct {
	repo *repository.OverviewRepository
}

func NewOverviewService(repo *repository.OverviewRepository) *OverviewService {
	return &OverviewService{repo: repo}
}

func (s *OverviewService) Get() (model.OverviewResponse, error) {
	return s.repo.GetSummary()
}
