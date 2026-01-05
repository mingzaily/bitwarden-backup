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
