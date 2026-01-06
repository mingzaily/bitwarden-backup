package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// GetServers 获取所有服务器配置
func GetServers(c *gin.Context) {
	var servers []database.ServerConfig
	query := database.DB

	if enabledParam := c.Query("enabled"); enabledParam != "" {
		enabled, err := strconv.ParseBool(enabledParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid enabled parameter"})
			return
		}
		query = query.Where("enabled = ?", enabled)
	}

	if err := query.Find(&servers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, servers)
}

// GetServer 获取单个服务器配置
func GetServer(c *gin.Context) {
	id := c.Param("id")
	var server database.ServerConfig
	if err := database.DB.First(&server, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}
	c.JSON(http.StatusOK, server)
}
