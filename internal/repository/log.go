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

// FindPaginated 分页查询日志，并带出所属任务名称。
func (r *LogRepository) FindPaginated(params model.PaginationParams, taskID *uint) ([]model.LogResponse, int64, error) {
	var logs []model.LogResponse
	var total int64

	query := r.db.Model(&model.BackupLog{})
	if taskID != nil {
		query = query.Where("task_id = ?", *taskID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	listQuery := r.db.Table("backup_logs AS logs").
		Select("logs.id, logs.task_id, COALESCE(tasks.name, '') AS task_name, logs.status, logs.message, logs.backup_file, logs.execution_logs, logs.start_time, logs.end_time, logs.created_at").
		Joins("LEFT JOIN backup_tasks AS tasks ON tasks.id = logs.task_id")
	if taskID != nil {
		listQuery = listQuery.Where("logs.task_id = ?", *taskID)
	}
	err := listQuery.
		Order("logs.created_at DESC").
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
