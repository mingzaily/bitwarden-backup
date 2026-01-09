package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// GetDestinations 获取所有备份目标（支持分页）
func GetDestinations(c *gin.Context) {
	var params model.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dests, total, err := destinationSvc.GetPaginated(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应结构
	responses := make([]model.DestinationResponse, len(dests))
	for i, d := range dests {
		responses[i] = d.ToResponse()
	}

	resp := model.NewPaginatedResponse(responses, params.Page, params.GetLimit(), total)
	c.JSON(http.StatusOK, resp)
}

// GetDestination 获取单个备份目标
func GetDestination(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	dest, err := destinationSvc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Destination not found"})
		return
	}
	c.JSON(http.StatusOK, dest.ToResponse())
}

// CreateDestination 创建备份目标
func CreateDestination(c *gin.Context) {
	var dest model.BackupDestination
	if err := c.ShouldBindJSON(&dest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := destinationSvc.Create(&dest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 重新获取以确保数据完整
	created, _ := destinationSvc.GetByID(dest.ID)
	c.JSON(http.StatusCreated, created.ToResponse())
}

// UpdateDestination 更新备份目标
func UpdateDestination(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req model.BackupDestination
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 判断是否仅更新 enabled 状态（前端 toggle 只传 enabled 字段）
	isToggleOnly := req.Name == "" && req.Type == ""

	if isToggleOnly {
		if err := destinationSvc.UpdateEnabled(uint(id), req.Enabled); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 获取现有记录
		existing, err := destinationSvc.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Destination not found"})
			return
		}

		// 更新非敏感字段
		existing.Name = req.Name
		existing.Type = req.Type
		existing.LocalPath = req.LocalPath
		existing.WebDAVURL = req.WebDAVURL
		existing.WebDAVUsername = req.WebDAVUsername
		existing.WebDAVPath = req.WebDAVPath
		existing.S3Endpoint = req.S3Endpoint
		existing.S3Region = req.S3Region
		existing.S3Bucket = req.S3Bucket
		existing.S3AccessKey = req.S3AccessKey
		existing.S3Path = req.S3Path
		existing.TargetServerID = req.TargetServerID
		existing.Encrypted = req.Encrypted
		existing.Enabled = req.Enabled
	// 更新备份保留份数配置
		if req.MaxBackupCount < 0 {
			req.MaxBackupCount = 0
		}
		existing.MaxBackupCount = req.MaxBackupCount


		// 敏感字段：空值不更新
		if req.WebDAVPassword != "" {
			existing.WebDAVPassword = req.WebDAVPassword
		}
		if req.S3SecretKey != "" {
			existing.S3SecretKey = req.S3SecretKey
		}
		if req.EncryptionPassword != "" {
			existing.EncryptionPassword = req.EncryptionPassword
		}

		if err := destinationSvc.Update(uint(id), existing); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 重新获取更新后的数据
	updated, _ := destinationSvc.GetByID(uint(id))
	c.JSON(http.StatusOK, updated.ToResponse())
}

// DeleteDestination 删除备份目标
func DeleteDestination(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := destinationSvc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Destination deleted"})
}

// ToggleDestination 切换备份目标启用状态
func ToggleDestination(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := destinationSvc.Toggle(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Destination toggled"})
}
