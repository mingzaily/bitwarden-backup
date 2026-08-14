package handler

import (
	"encoding/json"
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

type fakeLogService struct {
	deletedIDs []uint
	deleted    int64
	err        error
}

func (f *fakeLogService) GetPaginated(model.PaginationParams, *uint) ([]model.LogResponse, int64, error) {
	return nil, 0, nil
}

func (f *fakeLogService) DeleteByIDs(ids []uint) (int64, error) {
	f.deletedIDs = append([]uint(nil), ids...)
	return f.deleted, f.err
}

func TestDeleteLogsUsesInjectedService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeLogService{deleted: 2}
	api := NewWithDependencies(nil, nil, nil, service, nil)
	r := gin.New()
	r.DELETE("/logs", api.DeleteLogs)

	req := httptest.NewRequest(http.MethodDelete, "/logs", strings.NewReader(`{"ids":[7,9]}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", res.Code, http.StatusOK, res.Body.String())
	}
	if len(service.deletedIDs) != 2 || service.deletedIDs[0] != 7 || service.deletedIDs[1] != 9 {
		t.Fatalf("unexpected deleted IDs: %v", service.deletedIDs)
	}
	if !strings.Contains(res.Body.String(), `"deleted":2`) {
		t.Fatalf("unexpected response: %s", res.Body.String())
	}
}

func TestDeleteLogsRejectsInvalidSelections(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "empty", body: `{"ids":[]}`},
		{name: "zero ID", body: `{"ids":[0]}`},
		{name: "duplicate ID", body: `{"ids":[1,1]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			service := &fakeLogService{}
			api := NewWithDependencies(nil, nil, nil, service, nil)
			r := gin.New()
			r.DELETE("/logs", api.DeleteLogs)

			req := httptest.NewRequest(http.MethodDelete, "/logs", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body = %s", res.Code, http.StatusBadRequest, res.Body.String())
			}
			if service.deletedIDs != nil {
				t.Fatalf("service should not be called: %v", service.deletedIDs)
			}
		})
	}

	ids := make([]uint, maxLogDeleteBatchSize+1)
	for index := range ids {
		ids[index] = uint(index + 1)
	}
	body, err := json.Marshal(model.DeleteLogsRequest{IDs: ids})
	if err != nil {
		t.Fatalf("marshal oversized request: %v", err)
	}
	service := &fakeLogService{}
	api := NewWithDependencies(nil, nil, nil, service, nil)
	r := gin.New()
	r.DELETE("/logs", api.DeleteLogs)
	req := httptest.NewRequest(http.MethodDelete, "/logs", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("oversized status = %d, want %d; body = %s", res.Code, http.StatusBadRequest, res.Body.String())
	}
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
