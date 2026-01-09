package repository

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"gorm.io/gorm"
)

type DestinationRepository struct {
	db *gorm.DB
}

func NewDestinationRepository(db *gorm.DB) *DestinationRepository {
	return &DestinationRepository{db: db}
}

func (r *DestinationRepository) FindAll() ([]model.BackupDestination, error) {
	var dests []model.BackupDestination
	err := r.db.Preload("TargetServer").Find(&dests).Error
	return dests, err
}

func (r *DestinationRepository) FindByID(id uint) (*model.BackupDestination, error) {
	var dest model.BackupDestination
	err := r.db.Preload("TargetServer").First(&dest, id).Error
	return &dest, err
}

func (r *DestinationRepository) Create(dest *model.BackupDestination) error {
	return r.db.Create(dest).Error
}

func (r *DestinationRepository) Update(dest *model.BackupDestination) error {
	return r.db.Save(dest).Error
}

func (r *DestinationRepository) Delete(id uint) error {
	return r.db.Delete(&model.BackupDestination{}, id).Error
}
