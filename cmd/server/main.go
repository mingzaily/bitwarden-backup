package main

import (
	"log"

	"github.com/mingzaily/bitwarden-backup/internal/config"
	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/handler"
	"github.com/mingzaily/bitwarden-backup/internal/scheduler"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	if err := database.Init(cfg.DBPath, cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化 Handler 层
	handler.Init(database.DB)

	// 初始化调度器
	sched := scheduler.New()
	if err := sched.LoadTasks(); err != nil {
		log.Printf("Failed to load tasks: %v", err)
	}
	sched.Start()
	defer sched.Stop()

	// 初始化 Gin 路由
	r := setupRouter(cfg)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
