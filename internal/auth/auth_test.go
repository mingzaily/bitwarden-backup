package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoginSessionAndCSRF(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := New("correct horse battery staple", false)
	r := gin.New()
	r.POST("/login", m.Login)
	protected := r.Group("/protected")
	protected.Use(m.Require(), m.CSRF())
	protected.POST("", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	loginReq := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"password":"correct horse battery staple"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	r.ServeHTTP(loginRes, loginReq)
	if loginRes.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loginRes.Code, http.StatusOK)
	}
	var loginBody struct {
		CSRFToken string `json:"csrf_token"`
	}
	if err := json.Unmarshal(loginRes.Body.Bytes(), &loginBody); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if loginBody.CSRFToken == "" {
		t.Fatal("login did not return a CSRF token")
	}
	cookies := loginRes.Result().Cookies()
	if len(cookies) != 1 || !cookies[0].HttpOnly || cookies[0].Value == "" {
		t.Fatalf("unexpected session cookie: %#v", cookies)
	}

	withoutCSRF := httptest.NewRequest(http.MethodPost, "/protected", nil)
	withoutCSRF.AddCookie(cookies[0])
	withoutCSRFRes := httptest.NewRecorder()
	r.ServeHTTP(withoutCSRFRes, withoutCSRF)
	if withoutCSRFRes.Code != http.StatusForbidden {
		t.Fatalf("missing CSRF status = %d, want %d", withoutCSRFRes.Code, http.StatusForbidden)
	}

	protectedReq := httptest.NewRequest(http.MethodPost, "/protected", nil)
	protectedReq.AddCookie(cookies[0])
	protectedReq.Header.Set("X-CSRF-Token", loginBody.CSRFToken)
	protectedRes := httptest.NewRecorder()
	r.ServeHTTP(protectedRes, protectedReq)
	if protectedRes.Code != http.StatusNoContent {
		t.Fatalf("protected status = %d, want %d", protectedRes.Code, http.StatusNoContent)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := New("correct horse battery staple", false)
	r := gin.New()
	r.POST("/login", m.Login)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"password":"wrong password"}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("wrong password status = %d, want %d", res.Code, http.StatusUnauthorized)
	}
}
