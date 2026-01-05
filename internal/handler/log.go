package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// GetLogs 获取备份日志
func GetLogs(c *gin.Context) {
	taskID := c.Query("task_id")
	var logs []database.BackupLog

	query := database.DB.Order("created_at desc").Limit(100)
	if taskID != "" {
		query = query.Where("task_id = ?", taskID)
	}

	if err := query.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
