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

func (r *TaskRepository) Update(task *model.BackupTask) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) Delete(id uint) error {
	return r.db.Delete(&model.BackupTask{}, id).Error
}
