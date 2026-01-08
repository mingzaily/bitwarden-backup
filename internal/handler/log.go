package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetLogs 获取所有日志
func GetLogs(c *gin.Context) {
	logs, err := logSvc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
