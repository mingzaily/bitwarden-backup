package repository

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"gorm.io/gorm"
)

type ServerRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) FindAll() ([]model.ServerConfig, error) {
	var servers []model.ServerConfig
	err := r.db.Find(&servers).Error
	return servers, err
}

func (r *ServerRepository) FindByID(id uint) (*model.ServerConfig, error) {
	var server model.ServerConfig
	err := r.db.First(&server, id).Error
	return &server, err
}

func (r *ServerRepository) Create(server *model.ServerConfig) error {
	return r.db.Create(server).Error
}

func (r *ServerRepository) Update(server *model.ServerConfig) error {
	return r.db.Save(server).Error
}

func (r *ServerRepository) UpdateEnabled(id uint, enabled bool) error {
	return r.db.Model(&model.ServerConfig{}).Where("id = ?", id).Update("enabled", enabled).Error
}

func (r *ServerRepository) Delete(id uint) error {
	return r.db.Delete(&model.ServerConfig{}, id).Error
}
