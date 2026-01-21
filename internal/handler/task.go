package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// validateCronExpression 校验 Cron 表达式格式
func validateCronExpression(expr string) error {
	if expr == "" {
		return nil // 可选字段
	}
	parts := strings.Fields(expr)
	// 支持 5 位（分时日月周）或 6 位（秒分时日月周）
	if len(parts) != 5 && len(parts) != 6 {
		return errors.New("Cron 表达式格式不正确，应为 5 位或 6 位格式")
	}
	return nil
}

// validateSourceDestination 校验源和目标不能相同
func validateSourceDestination(sourceServerID uint, destinationIDs []uint) error {
	for _, destID := range destinationIDs {
		dest, err := destinationSvc.GetByID(destID)
		if err != nil {
			continue
		}
		// 如果目标类型是服务器，且目标服务器ID等于源服务器ID
		if dest.Type == "server" && dest.TargetServerID != nil && *dest.TargetServerID == sourceServerID {
			return errors.New("备份目标不能与源服务器相同")
		}
	}
	return nil
}

// GetTasks 获取所有任务（支持分页）
func GetTasks(c *gin.Context) {
	var params model.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks, total, err := taskSvc.GetPaginated(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
func GetTask(c *gin.Context) {
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
	c.JSON(http.StatusOK, task.ToResponse())
}

// CreateTask 创建任务
func CreateTask(c *gin.Context) {
	var req model.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 必填字段校验
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务名称不能为空"})
		return
	}
	if req.SourceServerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择源服务器"})
		return
	}
	if len(req.DestinationIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请至少选择一个备份目标"})
		return
	}

	// Cron 表达式格式校验
	if err := validateCronExpression(req.CronExpression); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 校验源和目标不能相同
	if err := validateSourceDestination(req.SourceServerID, req.DestinationIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &model.BackupTask{
		Name:           req.Name,
		SourceServerID: req.SourceServerID,
		CronExpression: req.CronExpression,
		Enabled:        true,
	}

	if err := taskSvc.CreateWithDestinations(task, req.DestinationIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回创建后的完整任务
	createdTask, _ := taskSvc.GetByID(task.ID)

	// 动态添加到调度器
	if taskScheduler != nil && createdTask.Enabled && createdTask.CronExpression != "" {
		if err := taskScheduler.AddTask(*createdTask); err != nil {
			// 记录错误但不影响创建结果
			c.JSON(http.StatusCreated, gin.H{
				"data":    createdTask.ToResponse(),
				"warning": "任务已创建，但添加到调度器失败: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, createdTask.ToResponse())
}

// UpdateTask 更新任务
func UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req model.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 判断是否仅更新 enabled 状态（前端 toggle 操作只传 enabled 字段）
	isToggleOnly := req.Enabled != nil && req.Name == "" && req.SourceServerID == 0 && len(req.DestinationIDs) == 0

	if isToggleOnly {
		// 仅更新 enabled 状态，不影响其他字段和关联
		if err := taskSvc.UpdateEnabled(uint(id), *req.Enabled); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 完整更新任务
		// Cron 表达式格式校验
		if err := validateCronExpression(req.CronExpression); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 校验源和目标不能相同
		if err := validateSourceDestination(req.SourceServerID, req.DestinationIDs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		task := &model.BackupTask{
			Name:           req.Name,
			SourceServerID: req.SourceServerID,
			CronExpression: req.CronExpression,
		}

		if req.Enabled != nil {
			task.Enabled = *req.Enabled
		}

		if err := taskSvc.UpdateWithDestinations(uint(id), task, req.DestinationIDs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 返回更新后的完整任务
	updatedTask, _ := taskSvc.GetByID(uint(id))

	// 动态更新调度器
	if taskScheduler != nil && updatedTask != nil {
		if err := taskScheduler.UpdateTask(*updatedTask); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"data":    updatedTask.ToResponse(),
				"warning": "任务已更新，但同步调度器失败: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, updatedTask.ToResponse())
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// 先从调度器移除任务
	if taskScheduler != nil {
		taskScheduler.RemoveTask(uint(id))
	}

	if err := taskSvc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
