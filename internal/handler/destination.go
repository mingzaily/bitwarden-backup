package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// GetDestinations 获取所有备份目标
func GetDestinations(c *gin.Context) {
	var destinations []database.BackupDestination
	if err := database.DB.Find(&destinations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, destinations)
}

// GetDestination 获取单个备份目标
func GetDestination(c *gin.Context) {
	id := c.Param("id")
	var destination database.BackupDestination
	if err := database.DB.First(&destination, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Destination not found"})
		return
	}
	c.JSON(http.StatusOK, destination)
}

// ToggleDestination 切换备份目标启用状态
func ToggleDestination(c *gin.Context) {
	id := c.Param("id")
	var destination database.BackupDestination
	if err := database.DB.First(&destination, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Destination not found"})
		return
	}

	var input struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	destination.Enabled = input.Enabled
	if err := database.DB.Save(&destination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, destination)
}
