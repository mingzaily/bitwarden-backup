package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/auth"
	"github.com/mingzaily/bitwarden-backup/internal/config"
	"github.com/mingzaily/bitwarden-backup/internal/handler"
)

func TestMetaRouteReturnsAppVersionWithoutAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRoutes(
		gin.New(),
		&config.Config{AppVersion: "v0.2.1"},
		auth.New("12345678", false),
		handler.NewWithDependencies(nil, nil, nil, nil, nil),
	)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/meta", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("GET /api/meta status = %d, want %d", res.Code, http.StatusOK)
	}

	var body struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode /api/meta response: %v", err)
	}
	if body.Version != "v0.2.1" {
		t.Fatalf("GET /api/meta version = %q, want v0.2.1", body.Version)
	}
}
