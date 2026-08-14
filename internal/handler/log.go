package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

const maxLogDeleteBatchSize = 100

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

// DeleteLogs removes selected completed execution records. Backup artifacts
// referenced by those records are intentionally left untouched.
func (a *API) DeleteLogs(c *gin.Context) {
	var request model.DeleteLogsRequest
	if !bindJSON(c, &request) {
		return
	}
	if len(request.IDs) == 0 {
		writeBadRequest(c, "请选择要删除的运行记录")
		return
	}
	if len(request.IDs) > maxLogDeleteBatchSize {
		writeBadRequest(c, "一次最多删除 100 条运行记录")
		return
	}
	seen := make(map[uint]struct{}, len(request.IDs))
	for _, id := range request.IDs {
		if id == 0 {
			writeBadRequest(c, "运行记录 ID 无效")
			return
		}
		if _, ok := seen[id]; ok {
			writeBadRequest(c, "运行记录不能重复选择")
			return
		}
		seen[id] = struct{}{}
	}

	deleted, err := a.logService.DeleteByIDs(request.IDs)
	if err != nil {
		writeInternalError(c, "delete logs", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": deleted})
}
