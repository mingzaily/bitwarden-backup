package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/scheduler"
)

// ExecuteTask 立即执行备份任务
func ExecuteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	task, err := taskSvc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 在后台执行任务
	go func() {
		sched := scheduler.New()
		sched.ExecuteTaskNow(*task)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Task execution started"})
}
