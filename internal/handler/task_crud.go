package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// CreateTask 创建备份任务
func CreateTask(c *gin.Context) {
	var input struct {
		Name             string `json:"name" binding:"required"`
		SourceServerID   uint   `json:"source_server_id" binding:"required"`
		CronExpression   string `json:"cron_expression"` // 可选，为空则仅支持手动触发
		Enabled          bool   `json:"enabled"`
		DestinationIDs   []uint `json:"destination_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := database.BackupTask{
		Name:           input.Name,
		SourceServerID: input.SourceServerID,
		CronExpression: input.CronExpression,
		Enabled:        input.Enabled,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 关联备份目标
	if len(input.DestinationIDs) > 0 {
		var destinations []database.BackupDestination
		database.DB.Find(&destinations, input.DestinationIDs)
		database.DB.Model(&task).Association("Destinations").Replace(destinations)
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTask 更新备份任务
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task database.BackupTask
	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input struct {
		Name             string `json:"name"`
		SourceServerID   uint   `json:"source_server_id"`
		CronExpression   string `json:"cron_expression"`
		Enabled          bool   `json:"enabled"`
		DestinationIDs   []uint `json:"destination_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.Name = input.Name
	task.SourceServerID = input.SourceServerID
	task.CronExpression = input.CronExpression
	task.Enabled = input.Enabled

	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新关联的备份目标
	var destinations []database.BackupDestination
	if len(input.DestinationIDs) > 0 {
		database.DB.Find(&destinations, input.DestinationIDs)
	}
	database.DB.Model(&task).Association("Destinations").Replace(destinations)

	c.JSON(http.StatusOK, task)
}

// DeleteTask 删除备份任务
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&database.BackupTask{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
