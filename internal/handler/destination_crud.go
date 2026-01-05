package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// CreateDestination 创建备份目标
func CreateDestination(c *gin.Context) {
	var destination database.BackupDestination
	if err := c.ShouldBindJSON(&destination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&destination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, destination)
}

// UpdateDestination 更新备份目标
func UpdateDestination(c *gin.Context) {
	id := c.Param("id")
	var destination database.BackupDestination
	if err := database.DB.First(&destination, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Destination not found"})
		return
	}

	if err := c.ShouldBindJSON(&destination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&destination).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, destination)
}

// DeleteDestination 删除备份目标
func DeleteDestination(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&database.BackupDestination{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Destination deleted"})
}
