package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// GetLogs 获取所有日志（支持分页）
func (a *API) GetLogs(c *gin.Context) {
	var params model.PaginationParams
	if !bindQuery(c, &params) {
		return
	}

	// 解析 task_id 参数
	taskID, ok := parseQueryID(c, "task_id")
	if !ok {
		return
	}

	// 分页查询
	logs, total, err := a.logService.GetPaginated(params, taskID)
	if err != nil {
		writeInternalError(c, "list logs", err)
		return
	}

	// 返回分页响应
	resp := model.NewPaginatedResponse(logs, params.Page, params.GetLimit(), total)
	c.JSON(http.StatusOK, resp)
}
