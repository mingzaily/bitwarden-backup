package repository

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"gorm.io/gorm"
)

type LogRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{db: db}
}

func (r *LogRepository) FindAll() ([]model.BackupLog, error) {
	var logs []model.BackupLog
	err := r.db.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *LogRepository) FindByTaskID(taskID uint) ([]model.BackupLog, error) {
	var logs []model.BackupLog
	err := r.db.Where("task_id = ?", taskID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *LogRepository) Create(log *model.BackupLog) error {
	return r.db.Create(log).Error
}

func (r *LogRepository) Update(log *model.BackupLog) error {
	return r.db.Save(log).Error
}
