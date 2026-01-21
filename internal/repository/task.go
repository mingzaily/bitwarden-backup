package repository

import (
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) FindAll() ([]model.BackupTask, error) {
	var tasks []model.BackupTask
	err := r.db.Preload("SourceServer").Preload("Destinations").Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) FindByID(id uint) (*model.BackupTask, error) {
	var task model.BackupTask
	err := r.db.Preload("SourceServer").Preload("Destinations").First(&task, id).Error
	return &task, err
}

func (r *TaskRepository) FindEnabled() ([]model.BackupTask, error) {
	var tasks []model.BackupTask
	err := r.db.Preload("SourceServer").Preload("Destinations").
		Where("enabled = ?", true).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) Create(task *model.BackupTask) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) CreateWithDestinations(task *model.BackupTask, destinationIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建任务
		if err := tx.Create(task).Error; err != nil {
			return err
		}

		// 添加目标关联
		if len(destinationIDs) > 0 {
			var destinations []model.BackupDestination
			if err := tx.Where("id IN ?", destinationIDs).Find(&destinations).Error; err != nil {
				return err
			}
			if err := tx.Model(task).Association("Destinations").Replace(destinations); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TaskRepository) Update(task *model.BackupTask) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) UpdateEnabled(id uint, enabled bool) error {
	return r.db.Model(&model.BackupTask{}).Where("id = ?", id).Update("enabled", enabled).Error
}

func (r *TaskRepository) UpdateWithDestinations(task *model.BackupTask, destinationIDs []uint) error {
	// 开启事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 只更新指定字段，保留 created_at
		if err := tx.Model(&model.BackupTask{}).Where("id = ?", task.ID).Updates(map[string]any{
			"name":             task.Name,
			"source_server_id": task.SourceServerID,
			"cron_expression":  task.CronExpression,
			"enabled":          task.Enabled,
		}).Error; err != nil {
			return err
		}

		// 清除旧的关联
		if err := tx.Model(task).Association("Destinations").Clear(); err != nil {
			return err
		}

		// 添加新的关联
		if len(destinationIDs) > 0 {
			var destinations []model.BackupDestination
			if err := tx.Where("id IN ?", destinationIDs).Find(&destinations).Error; err != nil {
				return err
			}
			if err := tx.Model(task).Association("Destinations").Replace(destinations); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TaskRepository) Delete(id uint) error {
	return r.db.Delete(&model.BackupTask{}, id).Error
}

// FindPaginated 分页查询任务
func (r *TaskRepository) FindPaginated(params model.PaginationParams) ([]model.BackupTask, int64, error) {
	var tasks []model.BackupTask
	var total int64

	if err := r.db.Model(&model.BackupTask{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("SourceServer").Preload("Destinations").
		Order("created_at DESC").
		Offset(params.GetOffset()).
		Limit(params.GetLimit()).
		Find(&tasks).Error

	return tasks, total, err
}
