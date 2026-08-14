package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOverview returns the aggregate data for the authenticated dashboard.
func (a *API) GetOverview(c *gin.Context) {
	if a.overviewService == nil {
		writeInternalError(c, "load overview", errors.New("overview service is not configured"))
		return
	}

	overview, err := a.overviewService.Get()
	if err != nil {
		writeInternalError(c, "load overview", err)
		return
	}
	c.JSON(http.StatusOK, overview)
}
