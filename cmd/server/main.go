package main

import (
	"log/slog"
	"os"

	"github.com/mingzaily/bitwarden-backup/internal/config"
	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/handler"
	"github.com/mingzaily/bitwarden-backup/internal/logger"
	"github.com/mingzaily/bitwarden-backup/internal/scheduler"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化日志
	logLevel := slog.LevelInfo
	if os.Getenv("BW_DEBUG") == "1" || os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}
	logger.Init(logLevel)

	// 初始化数据库
	if err := database.Init(cfg.DBPath, cfg); err != nil {
		logger.Module(logger.ModuleMain).Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	// 初始化 Handler 层
	handler.Init(database.DB)

	// 初始化调度器
	sched := scheduler.New()
	if err := sched.LoadTasks(); err != nil {
		logger.Module(logger.ModuleMain).Error("Failed to load tasks", "error", err)
	}
	sched.Start()
	defer sched.Stop()

	// 将调度器注入到 Handler 层，支持动态更新任务
	handler.SetScheduler(sched)

	// 初始化 Gin 路由
	r := setupRouter(cfg)

	// 启动服务器
	logger.Module(logger.ModuleMain).Info("Server starting", "port", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		logger.Module(logger.ModuleMain).Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
