package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// GetTasks 获取所有备份任务
func GetTasks(c *gin.Context) {
	var tasks []database.BackupTask
	// 预加载关联的备份目标
	if err := database.DB.Preload("Destinations").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTask 获取单个备份任务
func GetTask(c *gin.Context) {
	id := c.Param("id")
	var task database.BackupTask
	// 预加载关联的备份目标
	if err := database.DB.Preload("Destinations").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}
