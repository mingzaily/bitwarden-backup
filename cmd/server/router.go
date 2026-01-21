package main

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/config"
	"github.com/mingzaily/bitwarden-backup/internal/handler"
)

func setupRouter(cfg *config.Config) *gin.Engine {
	// 根据环境设置 Gin 模式
	if cfg.AppEnv == "dev" {
		gin.SetMode(gin.DebugMode)
		r := gin.Default() // 包含访问日志
		return setupRoutes(r)
	} else {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.Use(gin.Recovery()) // 仅保留 panic 恢复
		return setupRoutes(r)
	}
}

func setupRoutes(r *gin.Engine) *gin.Engine {

	// 静态资源（Vue 构建产物）
	r.Static("/assets", "./web/dist/assets")

	// API 路由
	api := r.Group("/api")
	{
		// 服务器配置
		api.GET("/servers", handler.GetServers)
		api.GET("/servers/:id", handler.GetServer)
		api.POST("/servers", handler.CreateServer)
		api.PUT("/servers/:id", handler.UpdateServer)
		api.DELETE("/servers/:id", handler.DeleteServer)

		// 备份目标
		api.GET("/destinations", handler.GetDestinations)
		api.GET("/destinations/:id", handler.GetDestination)
		api.POST("/destinations", handler.CreateDestination)
		api.PUT("/destinations/:id", handler.UpdateDestination)
		api.DELETE("/destinations/:id", handler.DeleteDestination)
		api.PATCH("/destinations/:id/toggle", handler.ToggleDestination)

		// 备份任务
		api.GET("/tasks", handler.GetTasks)
		api.GET("/tasks/:id", handler.GetTask)
		api.POST("/tasks", handler.CreateTask)
		api.PUT("/tasks/:id", handler.UpdateTask)
		api.DELETE("/tasks/:id", handler.DeleteTask)
		api.POST("/tasks/:id/execute", handler.ExecuteTask)

		// 日志
		api.GET("/logs", handler.GetLogs)
	}

	// SPA History Mode Fallback
	// 对于非 API 和非静态资源的请求，返回 index.html
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是 API 请求，返回 404
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// 如果是静态资源请求（有文件扩展名），返回 404
		if filepath.Ext(path) != "" {
			c.Status(http.StatusNotFound)
			return
		}

		// 其他请求返回 Vue SPA 的 index.html
		c.File("./web/dist/index.html")
	})

	return r
}
