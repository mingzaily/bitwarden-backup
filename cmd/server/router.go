package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/handler"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// 静态文件
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("./web/templates/*")

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

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

	return r
}
