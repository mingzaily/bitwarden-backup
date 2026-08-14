package repository

import (
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/model"
	"gorm.io/gorm"
)

const overviewRecentLimit = 5

// OverviewRepository contains the read-only aggregate queries used by the
// dashboard. Keeping these queries here avoids making the HTTP layer compose
// several resource services and prevents N+1 lookups for recent activity.
type OverviewRepository struct {
	db *gorm.DB
}

func NewOverviewRepository(db *gorm.DB) *OverviewRepository {
	return &OverviewRepository{db: db}
}

func (r *OverviewRepository) GetSummary() (model.OverviewResponse, error) {
	var response model.OverviewResponse

	if err := r.countResource(&response.Servers, &model.ServerConfig{}); err != nil {
		return response, err
	}
	if err := r.countResource(&response.Destinations, &model.BackupDestination{}); err != nil {
		return response, err
	}

	if err := r.db.Model(&model.BackupTask{}).Count(&response.Tasks.Total).Error; err != nil {
		return response, err
	}
	if err := r.db.Model(&model.BackupTask{}).
		Where("enabled = ?", true).
		Count(&response.Tasks.Enabled).Error; err != nil {
		return response, err
	}
	if err := r.db.Model(&model.BackupTask{}).
		Where("enabled = ? AND TRIM(cron_expression) <> ?", true, "").
		Count(&response.Tasks.Scheduled).Error; err != nil {
		return response, err
	}

	if err := r.db.Model(&model.BackupLog{}).Count(&response.Logs.Total).Error; err != nil {
		return response, err
	}
	cutoff := time.Now().Add(-24 * time.Hour)
	if err := r.countLogsSince(cutoff, "success", &response.Logs.Success24h); err != nil {
		return response, err
	}
	if err := r.countLogsSince(cutoff, "failed", &response.Logs.Failed24h); err != nil {
		return response, err
	}
	if err := r.countLogsSince(cutoff, "running", &response.Logs.Running24h); err != nil {
		return response, err
	}

	if err := r.loadRecentTasks(&response.RecentTasks); err != nil {
		return response, err
	}
	if err := r.loadRecentLogs(&response.RecentLogs); err != nil {
		return response, err
	}
	if response.RecentTasks == nil {
		response.RecentTasks = []model.OverviewTaskSummary{}
	}
	if response.RecentLogs == nil {
		response.RecentLogs = []model.OverviewLogSummary{}
	}

	return response, nil
}

func (r *OverviewRepository) countLogsSince(cutoff time.Time, status string, count *int64) error {
	return r.db.Model(&model.BackupLog{}).
		Where("created_at >= ? AND status = ?", cutoff, status).
		Count(count).Error
}

func (r *OverviewRepository) countResource(stats *model.OverviewCount, resource any) error {
	if err := r.db.Model(resource).Count(&stats.Total).Error; err != nil {
		return err
	}
	return r.db.Model(resource).Where("enabled = ?", true).Count(&stats.Enabled).Error
}

func (r *OverviewRepository) loadRecentTasks(tasks *[]model.OverviewTaskSummary) error {
	return r.db.Table("backup_tasks AS tasks").
		Select("tasks.id, tasks.name, tasks.enabled, tasks.cron_expression, COALESCE(servers.name, '') AS source_server_name, COUNT(task_destinations.backup_destination_id) AS destination_count, tasks.created_at").
		Joins("LEFT JOIN server_configs AS servers ON servers.id = tasks.source_server_id").
		Joins("LEFT JOIN task_destinations ON task_destinations.backup_task_id = tasks.id").
		Group("tasks.id, tasks.name, tasks.enabled, tasks.cron_expression, servers.name, tasks.created_at").
		Order("tasks.created_at DESC").
		Limit(overviewRecentLimit).
		Scan(tasks).Error
}

func (r *OverviewRepository) loadRecentLogs(logs *[]model.OverviewLogSummary) error {
	return r.db.Table("backup_logs AS logs").
		Select("logs.id, logs.task_id, COALESCE(tasks.name, '') AS task_name, logs.status, logs.message, logs.backup_file, logs.created_at").
		Joins("LEFT JOIN backup_tasks AS tasks ON tasks.id = logs.task_id").
		Order("logs.created_at DESC").
		Limit(overviewRecentLimit).
		Scan(logs).Error
}
