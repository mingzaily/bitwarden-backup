package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExecuteTask 立即执行备份任务
func (a *API) ExecuteTask(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	task, err := a.taskService.GetByID(id)
	if err != nil {
		writeLookupError(c, "task", "load task for execution", err)
		return
	}
	if !task.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task is disabled"})
		return
	}
	if a.scheduler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "scheduler is unavailable"})
		return
	}

	// 复用调度器队列，避免重复请求并发启动多个备份流程。
	if !a.scheduler.TriggerTask(task.ID) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "task is already running or scheduler queue is full"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Task execution queued"})
}
