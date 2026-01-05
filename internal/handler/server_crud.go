package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// CreateServer 创建服务器配置
func CreateServer(c *gin.Context) {
	var server database.ServerConfig
	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, server)
}

// UpdateServer 更新服务器配置
func UpdateServer(c *gin.Context) {
	id := c.Param("id")
	var server database.ServerConfig
	if err := database.DB.First(&server, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, server)
}

// DeleteServer 删除服务器配置
func DeleteServer(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&database.ServerConfig{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Server deleted"})
}
