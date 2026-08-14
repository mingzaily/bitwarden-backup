package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mingzaily/bitwarden-backup/internal/config"
	"github.com/mingzaily/bitwarden-backup/internal/crypto"
	applogger "github.com/mingzaily/bitwarden-backup/internal/logger"
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func Init(dbPath string, cfg *config.Config) error {
	if err := crypto.InitEncryption(); err != nil {
		return fmt.Errorf("failed to initialize encryption: %w", err)
	}

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create db directory: %w", err)
	}
	if dbPath != ":memory:" {
		file, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			return fmt.Errorf("failed to prepare database file: %w", err)
		}
		if err := file.Chmod(0600); err != nil {
			_ = file.Close()
			return fmt.Errorf("failed to secure database file: %w", err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close database file: %w", err)
		}
	}

	var logLevel logger.LogLevel
	if cfg.AppEnv == "dev" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Warn
	}

	// Critical #3: 启用外键约束，确保引用完整性
	// 注意：迁移时仍禁用外键以避免顺序问题，迁移完成后启用
	// SQLite is the control-plane database and is accessed by both HTTP
	// handlers and the scheduler. A busy timeout plus a single pooled
	// connection keeps short concurrent writes from surfacing as SQLITE_BUSY.
	dsn := dbPath + "?_pragma=foreign_keys%3d1&_pragma=busy_timeout%3d5000"
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dsn,
	}, &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true, // 迁移时暂时禁用，避免顺序问题
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to configure database pool: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	DB = db

	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	if err := MigrateEncryptExistingData(); err != nil {
		return fmt.Errorf("failed to encrypt existing data: %w", err)
	}

	applogger.Module(applogger.ModuleDatabase).Info("Database initialized successfully")
	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&model.ServerConfig{},
		&model.BackupTask{},
		&model.BackupDestination{},
		&model.BackupLog{},
	)
}
