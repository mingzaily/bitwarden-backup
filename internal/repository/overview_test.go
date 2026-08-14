package repository

import (
	"testing"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestOverviewRepositoryGetSummary(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:overview-repository-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(
		&model.ServerConfig{},
		&model.BackupTask{},
		&model.BackupDestination{},
		&model.BackupLog{},
	); err != nil {
		t.Fatalf("migrate database: %v", err)
	}

	server := model.ServerConfig{Name: "Production", ServerURL: "https://vault.example.com", Enabled: true}
	if err := db.Create(&server).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}
	disabledServer := model.ServerConfig{Name: "Disabled", ServerURL: "https://disabled.example.com"}
	if err := db.Create(&disabledServer).Error; err != nil {
		t.Fatalf("create disabled server: %v", err)
	}
	if err := db.Model(&disabledServer).Update("enabled", false).Error; err != nil {
		t.Fatalf("disable server: %v", err)
	}
	destination := model.BackupDestination{Name: "Local", Type: "local", LocalPath: "/backups", Enabled: true}
	if err := db.Create(&destination).Error; err != nil {
		t.Fatalf("create destination: %v", err)
	}
	task := model.BackupTask{Name: "Nightly", SourceServerID: server.ID, CronExpression: "0 0 2 * * *", Enabled: true}
	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("create task: %v", err)
	}
	if err := db.Model(&task).Association("Destinations").Append(&destination); err != nil {
		t.Fatalf("create task destination relation: %v", err)
	}
	if err := db.Create(&model.BackupLog{TaskID: task.ID, Status: "success", Message: "Backup completed successfully", CreatedAt: time.Now()}).Error; err != nil {
		t.Fatalf("create log: %v", err)
	}
	if err := db.Create(&model.BackupLog{TaskID: task.ID, Status: "failed", Message: "Backup failed", CreatedAt: time.Now()}).Error; err != nil {
		t.Fatalf("create failed log: %v", err)
	}

	overview, err := NewOverviewRepository(db).GetSummary()
	if err != nil {
		t.Fatalf("get overview: %v", err)
	}
	if overview.Servers.Total != 2 || overview.Servers.Enabled != 1 {
		t.Fatalf("unexpected server counts: %+v", overview.Servers)
	}
	if overview.Destinations.Total != 1 || overview.Destinations.Enabled != 1 {
		t.Fatalf("unexpected destination counts: %+v", overview.Destinations)
	}
	if overview.Tasks.Total != 1 || overview.Tasks.Enabled != 1 || overview.Tasks.Scheduled != 1 {
		t.Fatalf("unexpected task counts: %+v", overview.Tasks)
	}
	if overview.Logs.Total != 2 || overview.Logs.Success24h != 1 || overview.Logs.Failed24h != 1 {
		t.Fatalf("unexpected log counts: %+v", overview.Logs)
	}
	if len(overview.RecentTasks) != 1 || overview.RecentTasks[0].SourceServerName != "Production" || overview.RecentTasks[0].DestinationCount != 1 {
		t.Fatalf("unexpected recent tasks: %+v", overview.RecentTasks)
	}
	if len(overview.RecentLogs) != 2 || overview.RecentLogs[0].TaskName != "Nightly" {
		t.Fatalf("unexpected recent logs: %+v", overview.RecentLogs)
	}

	logs, total, err := NewLogRepository(db).FindPaginated(model.PaginationParams{Page: 1, PageSize: 10}, nil)
	if err != nil {
		t.Fatalf("get paginated logs: %v", err)
	}
	if total != 2 || len(logs) != 2 || logs[0].TaskName != "Nightly" {
		t.Fatalf("unexpected paginated logs: total=%d logs=%+v", total, logs)
	}
}
