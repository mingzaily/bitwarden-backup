package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLogRepositoryDeleteByIDsKeepsRunningLogs(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:log-repository-delete-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get database connection: %v", err)
	}
	defer sqlDB.Close()
	if err := db.AutoMigrate(&model.BackupLog{}); err != nil {
		t.Fatalf("migrate database: %v", err)
	}

	logs := []model.BackupLog{
		{TaskID: 1, Status: "success", Message: "success", CreatedAt: time.Now()},
		{TaskID: 1, Status: "failed", Message: "failed", CreatedAt: time.Now()},
		{TaskID: 1, Status: "running", Message: "running", CreatedAt: time.Now()},
	}
	if err := db.Create(&logs).Error; err != nil {
		t.Fatalf("create logs: %v", err)
	}

	deleted, err := NewLogRepository(db).DeleteByIDs([]uint{logs[0].ID, logs[1].ID, logs[2].ID})
	if err != nil {
		t.Fatalf("delete logs: %v", err)
	}
	if deleted != 2 {
		t.Fatalf("deleted = %d, want 2", deleted)
	}

	for _, id := range []uint{logs[0].ID, logs[1].ID} {
		var log model.BackupLog
		if err := db.First(&log, id).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatalf("log %d still exists or returned unexpected error: %v", id, err)
		}
	}
	var running model.BackupLog
	if err := db.First(&running, logs[2].ID).Error; err != nil {
		t.Fatalf("running log should remain: %v", err)
	}
}
