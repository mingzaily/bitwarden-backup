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

func (r *DestinationRepository) UpdateEnabled(id uint, enabled bool) error {
	return r.db.Model(&model.BackupDestination{}).Where("id = ?", id).Update("enabled", enabled).Error
}

func (r *DestinationRepository) Delete(id uint) error {
	return r.db.Delete(&model.BackupDestination{}, id).Error
}

// FindPaginated 分页查询备份目标
func (r *DestinationRepository) FindPaginated(params model.PaginationParams) ([]model.BackupDestination, int64, error) {
	var dests []model.BackupDestination
	var total int64

	if err := r.db.Model(&model.BackupDestination{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("TargetServer").
		Order("created_at DESC").
		Offset(params.GetOffset()).
		Limit(params.GetLimit()).
		Find(&dests).Error

	return dests, total, err
}
