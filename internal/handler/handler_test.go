package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

type fakeServerService struct {
	server       *model.ServerConfig
	updateErr    error
	updatedID    uint
	updatedValue bool
}

func (f *fakeServerService) GetByID(id uint) (*model.ServerConfig, error) {
	if f.server == nil {
		f.server = &model.ServerConfig{ID: id, Name: "test server", ServerURL: "https://vault.example.com"}
	}
	return f.server, nil
}

func (f *fakeServerService) GetPaginated(model.PaginationParams, *bool) ([]model.ServerConfig, int64, error) {
	return nil, 0, nil
}

func (f *fakeServerService) Create(*model.ServerConfig) error { return nil }

func (f *fakeServerService) Update(uint, *model.ServerConfig) error { return nil }

func (f *fakeServerService) UpdateEnabled(id uint, enabled bool) error {
	f.updatedID = id
	f.updatedValue = enabled
	return f.updateErr
}

func (f *fakeServerService) Delete(uint) error { return nil }

func TestSetServerEnabledUsesInjectedService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeServerService{}
	api := NewWithDependencies(service, nil, nil, nil, nil)
	r := gin.New()
	r.PATCH("/servers/:id/enabled", api.SetServerEnabled)

	req := httptest.NewRequest(http.MethodPatch, "/servers/7/enabled", strings.NewReader(`{"enabled":false}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", res.Code, http.StatusOK, res.Body.String())
	}
	if service.updatedID != 7 || service.updatedValue {
		t.Fatalf("unexpected update: id=%d enabled=%v", service.updatedID, service.updatedValue)
	}
}

func TestSetServerEnabledRequiresExplicitValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeServerService{}
	api := NewWithDependencies(service, nil, nil, nil, nil)
	r := gin.New()
	r.PATCH("/servers/:id/enabled", api.SetServerEnabled)

	req := httptest.NewRequest(http.MethodPatch, "/servers/7/enabled", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
	}
	if service.updatedID != 0 {
		t.Fatal("service should not be called when enabled is omitted")
	}
}

func TestSetServerEnabledDoesNotExposeInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeServerService{updateErr: errors.New("database credentials should stay private")}
	api := NewWithDependencies(service, nil, nil, nil, nil)
	r := gin.New()
	r.PATCH("/servers/:id/enabled", api.SetServerEnabled)

	req := httptest.NewRequest(http.MethodPatch, "/servers/7/enabled", strings.NewReader(`{"enabled":true}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusInternalServerError)
	}
	if strings.Contains(res.Body.String(), "database credentials") {
		t.Fatal("internal error leaked in response")
	}
}

type fakeOverviewService struct {
	response model.OverviewResponse
	err      error
}

func (f *fakeOverviewService) Get() (model.OverviewResponse, error) {
	return f.response, f.err
}

func TestGetOverviewUsesInjectedService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	api := NewWithDependencies(nil, nil, nil, nil, nil)
	api.SetOverviewService(&fakeOverviewService{response: model.OverviewResponse{
		Servers: model.OverviewCount{Total: 2, Enabled: 1},
	}})
	r := gin.New()
	r.GET("/overview", api.GetOverview)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/overview", nil))

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", res.Code, http.StatusOK, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), `"total":2`) || !strings.Contains(res.Body.String(), `"enabled":1`) {
		t.Fatalf("overview response did not contain injected data: %s", res.Body.String())
	}
}
