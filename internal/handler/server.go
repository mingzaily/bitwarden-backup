package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// GetServers 获取所有服务器配置（支持分页）
func (a *API) GetServers(c *gin.Context) {
	var params model.PaginationParams
	if !bindQuery(c, &params) {
		return
	}

	// 解析 enabled 参数
	enabled, ok := parseOptionalBoolQuery(c, "enabled")
	if !ok {
		return
	}

	servers, total, err := a.serverService.GetPaginated(params, enabled)
	if err != nil {
		writeInternalError(c, "list servers", err)
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
func (a *API) GetServer(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	server, err := a.serverService.GetByID(id)
	if err != nil {
		writeLookupError(c, "server", "load server", err)
		return
	}
	c.JSON(http.StatusOK, server.ToResponse())
}

// CreateServer 创建服务器配置
func (a *API) CreateServer(c *gin.Context) {
	var req model.ServerRequest
	if !bindJSON(c, &req) {
		return
	}

	if err := validateServerRequest(req, true); err != nil {
		writeBadRequest(c, err.Error())
		return
	}

	server := &model.ServerConfig{
		Name:           req.Name,
		ServerURL:      req.ServerURL,
		ClientID:       req.ClientID,
		ClientSecret:   req.ClientSecret,
		MasterPassword: req.MasterPassword,
		IsOfficial:     model.IsOfficialServerURL(req.ServerURL),
		Enabled:        true,
	}

	if err := a.serverService.Create(server); err != nil {
		writeInternalError(c, "create server", err)
		return
	}
	c.JSON(http.StatusCreated, server.ToResponse())
}

// UpdateServer 更新服务器配置
func (a *API) UpdateServer(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req model.ServerRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := validateServerRequest(req, false); err != nil {
		writeBadRequest(c, err.Error())
		return
	}

	// 获取现有记录并应用完整更新。
	existing, err := a.serverService.GetByID(id)
	if err != nil {
		writeLookupError(c, "server", "load server for update", err)
		return
	}

	existing.Name = req.Name
	existing.ServerURL = req.ServerURL
	if req.ClientID != "" {
		existing.ClientID = req.ClientID
	}
	existing.IsOfficial = model.IsOfficialServerURL(req.ServerURL)
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

	if err := a.serverService.Update(id, existing); err != nil {
		writeLookupError(c, "server", "update server", err)
		return
	}

	// 重新获取更新后的数据
	updated, err := a.serverService.GetByID(id)
	if err != nil {
		writeLookupError(c, "server", "load updated server", err)
		return
	}
	c.JSON(http.StatusOK, updated.ToResponse())
}

// SetServerEnabled updates only the server state. Keeping this separate from
// the full update endpoint avoids inferring intent from omitted fields.
func (a *API) SetServerEnabled(c *gin.Context) {
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
	if err := a.serverService.UpdateEnabled(id, *request.Enabled); err != nil {
		writeLookupError(c, "server", "update server status", err)
		return
	}

	server, err := a.serverService.GetByID(id)
	if err != nil {
		writeLookupError(c, "server", "load updated server", err)
		return
	}
	c.JSON(http.StatusOK, server.ToResponse())
}

// DeleteServer 删除服务器配置
func (a *API) DeleteServer(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := a.serverService.Delete(id); err != nil {
		writeLookupError(c, "server", "delete server", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Server deleted"})
}
