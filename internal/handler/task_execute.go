package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/scheduler"
)

// ExecuteTask 立即执行备份任务
func ExecuteTask(c *gin.Context) {
	id := c.Param("id")
	var task database.BackupTask

	// 预加载关联的备份目标
	if err := database.DB.Preload("Destinations").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 在后台执行任务
	go func() {
		sched := scheduler.New()
		sched.ExecuteTaskNow(task)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Task execution started"})
}
