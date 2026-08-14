package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// validateCronExpression 校验 Cron 表达式格式
// validateSourceDestination 校验源和目标不能相同
func (a *API) validateSourceDestination(sourceServerID uint, destinationIDs []uint) error {
	if len(destinationIDs) > 100 {
		return errors.New("备份目标数量不能超过 100 个")
	}
	if _, err := a.serverService.GetByID(sourceServerID); err != nil {
		return errors.New("源服务器不存在")
	}
	seen := make(map[uint]struct{}, len(destinationIDs))
	for _, destID := range destinationIDs {
		if _, ok := seen[destID]; ok {
			return errors.New("备份目标不能重复")
		}
		seen[destID] = struct{}{}
		dest, err := a.destinationService.GetByID(destID)
		if err != nil {
			return errors.New("备份目标不存在")
		}
		// 如果目标类型是服务器，且目标服务器ID等于源服务器ID
		if dest.Type == "server" && dest.TargetServerID != nil && *dest.TargetServerID == sourceServerID {
			return errors.New("备份目标不能与源服务器相同")
		}
	}
	return nil
}

// GetTasks 获取所有任务（支持分页）
func (a *API) GetTasks(c *gin.Context) {
	var params model.PaginationParams
	if !bindQuery(c, &params) {
		return
	}

	tasks, total, err := a.taskService.GetPaginated(params)
	if err != nil {
		writeInternalError(c, "list tasks", err)
		return
	}

	// 转换为响应结构
	responses := make([]model.TaskResponse, len(tasks))
	for i, t := range tasks {
		responses[i] = t.ToResponse()
	}

	resp := model.NewPaginatedResponse(responses, params.Page, params.GetLimit(), total)
	c.JSON(http.StatusOK, resp)
}

// GetTask 获取单个任务
func (a *API) GetTask(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	task, err := a.taskService.GetByID(id)
	if err != nil {
		writeLookupError(c, "task", "load task", err)
		return
	}
	c.JSON(http.StatusOK, task.ToResponse())
}

// CreateTask 创建任务
func (a *API) CreateTask(c *gin.Context) {
	var req model.TaskRequest
	if !bindJSON(c, &req) {
		return
	}

	// 必填字段校验
	if req.Name == "" {
		writeBadRequest(c, "任务名称不能为空")
		return
	}
	if err := validateText(req.Name, "name", 100, true); err != nil {
		writeBadRequest(c, err.Error())
		return
	}
	if err := model.ValidateFilenameTemplate(req.FilenameTemplate); err != nil {
		writeBadRequest(c, err.Error())
		return
	}
	if req.SourceServerID == 0 {
		writeBadRequest(c, "请选择源服务器")
		return
	}
	if len(req.DestinationIDs) == 0 {
		writeBadRequest(c, "请至少选择一个备份目标")
		return
	}

	// Cron 表达式格式校验
	if err := validateCronExpression(req.CronExpression); err != nil {
		writeBadRequest(c, err.Error())
		return
	}

	// 校验源和目标不能相同
	if err := a.validateSourceDestination(req.SourceServerID, req.DestinationIDs); err != nil {
		writeBadRequest(c, err.Error())
		return
	}

	task := &model.BackupTask{
		Name:             req.Name,
		SourceServerID:   req.SourceServerID,
		CronExpression:   req.CronExpression,
		FilenameTemplate: model.NormalizeFilenameTemplate(req.FilenameTemplate),
		Enabled:          true,
	}

	if err := a.taskService.CreateWithDestinations(task, req.DestinationIDs); err != nil {
		writeInternalError(c, "create task", err)
		return
	}

	// 返回创建后的完整任务
	createdTask, err := a.taskService.GetByID(task.ID)
	if err != nil {
		writeLookupError(c, "task", "load created task", err)
		return
	}

	// 动态添加到调度器
	if a.scheduler != nil && createdTask.Enabled && createdTask.CronExpression != "" {
		if err := a.scheduler.AddTask(*createdTask); err != nil {
			// 记录错误但不影响创建结果
			c.JSON(http.StatusCreated, gin.H{
				"data":    createdTask.ToResponse(),
				"warning": "任务已创建，但添加到调度器失败",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, createdTask.ToResponse())
}

// UpdateTask 更新任务
func (a *API) UpdateTask(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req model.TaskRequest
	if !bindJSON(c, &req) {
		return
	}

	if err := validateText(req.Name, "name", 100, true); err != nil {
		writeBadRequest(c, err.Error())
		return
	}
	if err := model.ValidateFilenameTemplate(req.FilenameTemplate); err != nil {
		writeBadRequest(c, err.Error())
		return
	}
	if req.SourceServerID == 0 || len(req.DestinationIDs) == 0 {
		writeBadRequest(c, "源服务器和备份目标不能为空")
		return
	}
	if err := validateCronExpression(req.CronExpression); err != nil {
		writeBadRequest(c, err.Error())
		return
	}
	if err := a.validateSourceDestination(req.SourceServerID, req.DestinationIDs); err != nil {
		writeBadRequest(c, err.Error())
		return
	}

	task := &model.BackupTask{
		Name:             req.Name,
		SourceServerID:   req.SourceServerID,
		CronExpression:   req.CronExpression,
		FilenameTemplate: model.NormalizeFilenameTemplate(req.FilenameTemplate),
	}

	existing, err := a.taskService.GetByID(id)
	if err != nil {
		writeLookupError(c, "task", "load task for update", err)
		return
	}
	if req.Enabled != nil {
		task.Enabled = *req.Enabled
	} else {
		task.Enabled = existing.Enabled
	}
	// Keep a custom template when an older API client omits the newly added
	// field during an update. New tasks still receive the default above.
	if strings.TrimSpace(req.FilenameTemplate) == "" {
		task.FilenameTemplate = model.NormalizeFilenameTemplate(existing.FilenameTemplate)
	}

	if err := a.taskService.UpdateWithDestinations(id, task, req.DestinationIDs); err != nil {
		writeLookupError(c, "task", "update task", err)
		return
	}

	// 返回更新后的完整任务
	updatedTask, err := a.taskService.GetByID(id)
	if err != nil {
		writeLookupError(c, "task", "load updated task", err)
		return
	}

	// 动态更新调度器
	if a.scheduler != nil && updatedTask != nil {
		if err := a.scheduler.UpdateTask(*updatedTask); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"data":    updatedTask.ToResponse(),
				"warning": "任务已更新，但同步调度器失败",
			})
			return
		}
	}

	c.JSON(http.StatusOK, updatedTask.ToResponse())
}

// SetTaskEnabled updates only the task state. It is intentionally separate
// from the full task update endpoint.
func (a *API) SetTaskEnabled(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var request model.EnabledRequest
	if !bindJSON(c, &request) {
		return
	}
	if request.Enabled == nil {
		writeBadRequest(c, "enabled is required")
		return
	}
	if err := a.taskService.UpdateEnabled(id, *request.Enabled); err != nil {
		writeLookupError(c, "task", "update task status", err)
		return
	}

	task, err := a.taskService.GetByID(id)
	if err != nil {
		writeLookupError(c, "task", "load updated task", err)
		return
	}
	if a.scheduler != nil {
		if err := a.scheduler.UpdateTask(*task); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"data":    task.ToResponse(),
				"warning": "任务状态已更新，但同步调度器失败",
			})
			return
		}
	}
	c.JSON(http.StatusOK, task.ToResponse())
}

// DeleteTask 删除任务
func (a *API) DeleteTask(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	// 先从调度器移除任务
	if a.scheduler != nil {
		a.scheduler.RemoveTask(id)
	}

	if err := a.taskService.Delete(id); err != nil {
		writeLookupError(c, "task", "delete task", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
