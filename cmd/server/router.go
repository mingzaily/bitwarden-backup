package main

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/auth"
	"github.com/mingzaily/bitwarden-backup/internal/config"
	"github.com/mingzaily/bitwarden-backup/internal/handler"
)

func setupRouter(cfg *config.Config, authManager *auth.Manager, api *handler.API) *gin.Engine {
	// 根据环境设置 Gin 模式
	if cfg.AppEnv == "dev" {
		gin.SetMode(gin.DebugMode)
		r := gin.Default() // 包含访问日志
		return setupRoutes(r, cfg, authManager, api)
	} else {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.Use(gin.Recovery()) // 仅保留 panic 恢复
		return setupRoutes(r, cfg, authManager, api)
	}
}

func setupRoutes(r *gin.Engine, cfg *config.Config, authManager *auth.Manager, apiHandler *handler.API) *gin.Engine {
	r.Use(securityHeaders())
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 静态资源（Vue 构建产物）
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/favicon.svg", "./web/dist/favicon.svg")

	// API 路由
	api := r.Group("/api")
	api.Use(requestBodyLimit(1 << 20))
	{
		// 登录和版本信息保持公开，其余 API 都需要管理员会话。
		api.POST("/auth/login", authManager.Login)
		api.GET("/meta", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"version": cfg.AppVersion})
		})
		protected := api.Group("")
		protected.Use(authManager.Require(), authManager.CSRF())
		protected.GET("/auth/session", authManager.Session)
		protected.POST("/auth/logout", authManager.Logout)
		protected.GET("/overview", apiHandler.GetOverview)

		// 服务器配置
		protected.GET("/servers", apiHandler.GetServers)
		protected.GET("/servers/:id", apiHandler.GetServer)
		protected.POST("/servers", apiHandler.CreateServer)
		protected.PUT("/servers/:id", apiHandler.UpdateServer)
		protected.PATCH("/servers/:id/enabled", apiHandler.SetServerEnabled)
		protected.DELETE("/servers/:id", apiHandler.DeleteServer)

		// 备份目标
		protected.GET("/destinations", apiHandler.GetDestinations)
		protected.GET("/destinations/:id", apiHandler.GetDestination)
		protected.POST("/destinations", apiHandler.CreateDestination)
		protected.PUT("/destinations/:id", apiHandler.UpdateDestination)
		protected.POST("/destinations/:id/test", apiHandler.TestDestination)
		protected.PATCH("/destinations/:id/enabled", apiHandler.SetDestinationEnabled)
		protected.DELETE("/destinations/:id", apiHandler.DeleteDestination)
		protected.PATCH("/destinations/:id/toggle", apiHandler.ToggleDestination)

		// 备份任务
		protected.GET("/tasks", apiHandler.GetTasks)
		protected.GET("/tasks/:id", apiHandler.GetTask)
		protected.POST("/tasks", apiHandler.CreateTask)
		protected.PUT("/tasks/:id", apiHandler.UpdateTask)
		protected.PATCH("/tasks/:id/enabled", apiHandler.SetTaskEnabled)
		protected.DELETE("/tasks/:id", apiHandler.DeleteTask)
		protected.POST("/tasks/:id/execute", apiHandler.ExecuteTask)

		// 日志
		protected.GET("/logs", apiHandler.GetLogs)
		protected.DELETE("/logs", apiHandler.DeleteLogs)
	}

	// SPA History Mode Fallback
	// 对于非 API 和非静态资源的请求，返回 index.html
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是 API 请求，返回 404
		if path == "/api" || strings.HasPrefix(path, "/api/") {
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

func requestBodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		if strings.HasPrefix(c.Request.URL.Path, "/api/") || c.Request.URL.Path == "/api" {
			c.Header("Cache-Control", "no-store")
		}
		c.Header("Content-Security-Policy", "default-src 'self'; base-uri 'self'; frame-ancestors 'none'; object-src 'none'; img-src 'self' data:; style-src 'self'; script-src 'self'; connect-src 'self'")
		c.Next()
	}
}
