package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// GetDestinations 获取所有备份目标（支持分页）
func (a *API) GetDestinations(c *gin.Context) {
	var params model.PaginationParams
	if !bindQuery(c, &params) {
		return
	}

	destinations, total, err := a.destinationService.GetPaginated(params)
	if err != nil {
		writeInternalError(c, "list destinations", err)
		return
	}

	responses := make([]model.DestinationResponse, len(destinations))
	for i, destination := range destinations {
		responses[i] = destination.ToResponse()
	}

	response := model.NewPaginatedResponse(responses, params.Page, params.GetLimit(), total)
	c.JSON(http.StatusOK, response)
}

// GetDestination 获取单个备份目标
func (a *API) GetDestination(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	destination, err := a.destinationService.GetByID(id)
	if err != nil {
		writeLookupError(c, "destination", "load destination", err)
		return
	}
	c.JSON(http.StatusOK, destination.ToResponse())
}

// CreateDestination 创建备份目标
func (a *API) CreateDestination(c *gin.Context) {
	var request model.DestinationRequest
	if !bindJSON(c, &request) {
		return
	}
	if err := validateDestination(request); err != nil {
		writeBadRequest(c, err.Error())
		return
	}

	destination := request.ToDestination()
	if destination.Type == "server" && !a.validateTargetServer(c, destination.TargetServerID, "target server") {
		return
	}
	if err := a.destinationService.Create(destination); err != nil {
		writeInternalError(c, "create destination", err)
		return
	}

	created, err := a.destinationService.GetByID(destination.ID)
	if err != nil {
		writeLookupError(c, "destination", "load created destination", err)
		return
	}
	c.JSON(http.StatusCreated, created.ToResponse())
}

// UpdateDestination 更新备份目标
func (a *API) UpdateDestination(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var request model.DestinationRequest
	if !bindJSON(c, &request) {
		return
	}
	if err := validateDestination(request); err != nil {
		writeBadRequest(c, err.Error())
		return
	}

	destination, err := a.destinationService.GetByID(id)
	if err != nil {
		writeLookupError(c, "destination", "load destination for update", err)
		return
	}

	// The UI sends a masked access key when the key was not changed. Treat it
	// as omitted instead of persisting the mask as a real credential.
	if strings.Contains(request.S3AccessKey, "****") {
		request.S3AccessKey = ""
	}
	request.ApplyTo(destination)
	if destination.Type == "server" && !a.validateTargetServer(c, destination.TargetServerID, "target server") {
		return
	}

	if err := a.destinationService.Update(id, destination); err != nil {
		writeLookupError(c, "destination", "update destination", err)
		return
	}

	updated, err := a.destinationService.GetByID(id)
	if err != nil {
		writeLookupError(c, "destination", "load updated destination", err)
		return
	}
	c.JSON(http.StatusOK, updated.ToResponse())
}

// SetDestinationEnabled updates the destination state without overloading the
// full destination update contract.
func (a *API) SetDestinationEnabled(c *gin.Context) {
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
	if err := a.destinationService.UpdateEnabled(id, *request.Enabled); err != nil {
		writeLookupError(c, "destination", "update destination status", err)
		return
	}

	destination, err := a.destinationService.GetByID(id)
	if err != nil {
		writeLookupError(c, "destination", "load updated destination", err)
		return
	}
	c.JSON(http.StatusOK, destination.ToResponse())
}

// DeleteDestination 删除备份目标
func (a *API) DeleteDestination(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := a.destinationService.Delete(id); err != nil {
		writeLookupError(c, "destination", "delete destination", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Destination deleted"})
}

// ToggleDestination is kept for clients using the original toggle endpoint.
func (a *API) ToggleDestination(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := a.destinationService.Toggle(id); err != nil {
		writeLookupError(c, "destination", "toggle destination", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Destination toggled"})
}

func (a *API) validateTargetServer(c *gin.Context, serverID *uint, resource string) bool {
	if serverID == nil || *serverID == 0 {
		writeBadRequest(c, resource+" is required")
		return false
	}
	if _, err := a.serverService.GetByID(*serverID); err != nil {
		if isRecordNotFound(err) {
			writeBadRequest(c, resource+" not found")
		} else {
			writeInternalError(c, "load "+resource, err)
		}
		return false
	}
	return true
}
