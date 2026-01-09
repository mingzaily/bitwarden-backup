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

// FindPaginated 分页查询日志
func (r *LogRepository) FindPaginated(params model.PaginationParams, taskID *uint) ([]model.BackupLog, int64, error) {
	var logs []model.BackupLog
	var total int64

	query := r.db.Model(&model.BackupLog{})
	if taskID != nil {
		query = query.Where("task_id = ?", *taskID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").
		Offset(params.GetOffset()).
		Limit(params.GetLimit()).
		Find(&logs).Error

	return logs, total, err
}

func (r *LogRepository) Create(log *model.BackupLog) error {
	return r.db.Create(log).Error
}

func (r *LogRepository) Update(log *model.BackupLog) error {
	return r.db.Save(log).Error
}
