package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mingzaily/bitwarden-backup/internal/config"
	"github.com/mingzaily/bitwarden-backup/internal/crypto"
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
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create db directory: %w", err)
	}

	var logLevel logger.LogLevel
	if cfg.AppEnv == "dev" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Warn
	}

	dsn := dbPath + "?_pragma=foreign_keys(0)"
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dsn,
	}, &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	DB = db

	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database initialized successfully")
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
