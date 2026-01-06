package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mingzaily/bitwarden-backup/internal/crypto"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite" // 纯 Go SQLite 驱动
)

var DB *gorm.DB

// Init 初始化数据库连接
func Init(dbPath string) error {
	// 初始化加密系统
	if err := crypto.InitEncryption(); err != nil {
		return fmt.Errorf("failed to initialize encryption: %w", err)
	}

	// 确保数据库目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create db directory: %w", err)
	}

	// 使用纯 Go SQLite 驱动（不需要 CGO）
	dsn := dbPath + "?_pragma=foreign_keys(1)"
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dsn,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	DB = db

	// 自动迁移数据库表
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// autoMigrate 自动迁移数据库表结构
func autoMigrate() error {
	return DB.AutoMigrate(
		&ServerConfig{},
		&BackupTask{},
		&BackupDestination{},
		&BackupLog{},
	)
}
