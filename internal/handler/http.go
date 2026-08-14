package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/logger"
	"gorm.io/gorm"
)

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 32)
	if err != nil || id == 0 {
		writeBadRequest(c, "invalid id")
		return 0, false
	}
	return uint(id), true
}

func parseQueryID(c *gin.Context, key string) (*uint, bool) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil, true
	}

	id, err := strconv.ParseUint(raw, 10, 32)
	if err != nil || id == 0 {
		writeBadRequest(c, "invalid "+key)
		return nil, false
	}
	value := uint(id)
	return &value, true
}

func parseOptionalBoolQuery(c *gin.Context, key string) (*bool, bool) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil, true
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		writeBadRequest(c, "invalid "+key)
		return nil, false
	}
	return &value, true
}

func bindJSON(c *gin.Context, destination any) bool {
	if err := c.ShouldBindJSON(destination); err != nil {
		writeBadRequest(c, "invalid request body")
		return false
	}
	return true
}

func bindQuery(c *gin.Context, destination any) bool {
	if err := c.ShouldBindQuery(destination); err != nil {
		writeBadRequest(c, "invalid query parameters")
		return false
	}
	return true
}

func writeBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func writeNotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, gin.H{"error": resource + " not found"})
}

func writeInternalError(c *gin.Context, operation string, err error) {
	logger.Module(logger.ModuleHandler).Error(operation, "error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

func isRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func writeLookupError(c *gin.Context, resource, operation string, err error) {
	if isRecordNotFound(err) {
		writeNotFound(c, resource)
		return
	}
	writeInternalError(c, operation, err)
}
