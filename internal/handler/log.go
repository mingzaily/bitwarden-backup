package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetLogs 获取所有日志
func GetLogs(c *gin.Context) {
	taskIDStr := c.Query("task_id")

	var logs any
	var err error

	if taskIDStr != "" {
		taskID, parseErr := strconv.ParseUint(taskIDStr, 10, 32)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task_id"})
			return
		}
		logs, err = logSvc.GetByTaskID(uint(taskID))
	} else {
		logs, err = logSvc.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
