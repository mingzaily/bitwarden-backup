package service

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/repository"
)

type ServerService struct {
	repo *repository.ServerRepository
}

func NewServerService(repo *repository.ServerRepository) *ServerService {
	return &ServerService{repo: repo}
}

func (s *ServerService) GetAll() ([]model.ServerConfig, error) {
	return s.repo.FindAll()
}

func (s *ServerService) GetByID(id uint) (*model.ServerConfig, error) {
	return s.repo.FindByID(id)
}

func (s *ServerService) Create(server *model.ServerConfig) error {
	return s.repo.Create(server)
}

func (s *ServerService) Update(id uint, server *model.ServerConfig) error {
	server.ID = id
	return s.repo.Update(server)
}

func (s *ServerService) UpdateEnabled(id uint, enabled bool) error {
	return s.repo.UpdateEnabled(id, enabled)
}

func (s *ServerService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetPaginated 分页获取服务器
func (s *ServerService) GetPaginated(params model.PaginationParams, enabled *bool) ([]model.ServerConfig, int64, error) {
	return s.repo.FindPaginated(params, enabled)
}
