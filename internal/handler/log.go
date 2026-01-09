package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// GetLogs 获取所有日志（支持分页）
func GetLogs(c *gin.Context) {
	var params model.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析 task_id 参数
	var taskID *uint
	if taskIDStr := c.Query("task_id"); taskIDStr != "" {
		id, err := strconv.ParseUint(taskIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task_id"})
			return
		}
		tid := uint(id)
		taskID = &tid
	}

	// 分页查询
	logs, total, err := logSvc.GetPaginated(params, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回分页响应
	resp := model.NewPaginatedResponse(logs, params.Page, params.GetLimit(), total)
	c.JSON(http.StatusOK, resp)
}
