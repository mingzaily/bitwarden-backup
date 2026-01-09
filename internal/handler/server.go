package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// GetServers 获取所有服务器配置（支持分页）
func GetServers(c *gin.Context) {
	var params model.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析 enabled 参数
	var enabled *bool
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		e := enabledStr == "true"
		enabled = &e
	}

	servers, total, err := serverSvc.GetPaginated(params, enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 转换为响应结构
	responses := make([]model.ServerResponse, len(servers))
	for i, s := range servers {
		responses[i] = s.ToResponse()
	}

	resp := model.NewPaginatedResponse(responses, params.Page, params.GetLimit(), total)
	c.JSON(http.StatusOK, resp)
}

// GetServer 获取单个服务器配置
func GetServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	server, err := serverSvc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}
	c.JSON(http.StatusOK, server.ToResponse())
}

// CreateServer 创建服务器配置
func CreateServer(c *gin.Context) {
	var req model.ServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 新增时敏感字段不能为空
	if req.ClientSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_secret is required"})
		return
	}
	if req.MasterPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "master_password is required"})
		return
	}

	server := &model.ServerConfig{
		Name:           req.Name,
		ServerURL:      req.ServerURL,
		ClientID:       req.ClientID,
		ClientSecret:   req.ClientSecret,
		MasterPassword: req.MasterPassword,
		IsOfficial:     req.IsOfficial,
		Enabled:        true,
	}

	if err := serverSvc.Create(server); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, server.ToResponse())
}

// UpdateServer 更新服务器配置
func UpdateServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req model.ServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 判断是否仅更新 enabled 状态
	isToggleOnly := req.Enabled != nil && req.Name == "" && req.ServerURL == "" && req.ClientID == ""

	if isToggleOnly {
		if err := serverSvc.UpdateEnabled(uint(id), *req.Enabled); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 获取现有记录
		existing, err := serverSvc.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
			return
		}

		// 更新非敏感字段
		existing.Name = req.Name
		existing.ServerURL = req.ServerURL
		existing.ClientID = req.ClientID
		existing.IsOfficial = req.IsOfficial
		if req.Enabled != nil {
			existing.Enabled = *req.Enabled
		}

		// 敏感字段：空值不更新
		if req.ClientSecret != "" {
			existing.ClientSecret = req.ClientSecret
		}
		if req.MasterPassword != "" {
			existing.MasterPassword = req.MasterPassword
		}

		if err := serverSvc.Update(uint(id), existing); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 重新获取更新后的数据
	updated, _ := serverSvc.GetByID(uint(id))
	c.JSON(http.StatusOK, updated.ToResponse())
}

// DeleteServer 删除服务器配置
func DeleteServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := serverSvc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Server deleted"})
}
