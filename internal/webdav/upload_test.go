package webdav

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRequestURLJoinsBaseAndRemotePath(t *testing.T) {
	client := NewClient("https://dav.example.test/root/", "user", "password")

	got, err := client.requestURL("/vault/backup.zip")
	if err != nil {
		t.Fatalf("requestURL returned error: %v", err)
	}
	if want := "https://dav.example.test/root/vault/backup.zip"; got != want {
		t.Fatalf("requestURL() = %q, want %q", got, want)
	}
}

func TestRequestURLRejectsParentTraversal(t *testing.T) {
	client := NewClient("https://dav.example.test/root", "user", "password")

	for _, remotePath := range []string{"../backup.zip", "/vault/../backup.zip", "vault\\backup.zip"} {
		if _, err := client.requestURL(remotePath); err == nil {
			t.Errorf("requestURL(%q) accepted an unsafe path", remotePath)
		}
	}
}

func TestRequestURLRejectsInsecureRemoteBase(t *testing.T) {
	client := NewClient("http://dav.example.test/root", "user", "password")
	if _, err := client.requestURL("backup.zip"); err == nil {
		t.Fatal("requestURL accepted a non-loopback HTTP base URL")
	}
}

func TestUploadFileCreatesMissingRemoteDirectories(t *testing.T) {
	tempDir := t.TempDir()
	localPath := filepath.Join(tempDir, "backup.json")
	wantBody := `{"backup":true}`
	if err := os.WriteFile(localPath, []byte(wantBody), 0600); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	collections := map[string]bool{"/dav/": true}
	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, password, ok := r.BasicAuth(); !ok || user != "user" || password != "password" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		requests = append(requests, r.Method+" "+r.URL.Path)

		switch r.Method {
		case "PROPFIND":
			if !collections[r.URL.Path] {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusMultiStatus)
		case "MKCOL":
			parent := path.Dir(strings.TrimSuffix(r.URL.Path, "/")) + "/"
			if !collections[parent] {
				w.WriteHeader(http.StatusConflict)
				return
			}
			collections[r.URL.Path] = true
			w.WriteHeader(http.StatusCreated)
		case "PUT":
			parent := path.Dir(r.URL.Path) + "/"
			if !collections[parent] {
				w.WriteHeader(http.StatusConflict)
				return
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("failed to read uploaded body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if string(body) != wantBody {
				t.Errorf("uploaded body = %q, want %q", body, wantBody)
			}
			if r.ContentLength != int64(len(wantBody)) {
				t.Errorf("Content-Length = %d, want %d", r.ContentLength, len(wantBody))
			}
			if got := r.Header.Get("Content-Type"); got != "application/octet-stream" {
				t.Errorf("Content-Type = %q, want application/octet-stream", got)
			}
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL+"/dav/", "user", "password")
	if err := client.UploadFile(localPath, "/nested/deep/backup.json"); err != nil {
		t.Fatalf("UploadFile returned error: %v", err)
	}

	wantRequests := []string{
		"PROPFIND /dav/nested/",
		"MKCOL /dav/nested/",
		"PROPFIND /dav/nested/deep/",
		"MKCOL /dav/nested/deep/",
		"PUT /dav/nested/deep/backup.json",
	}
	if !reflect.DeepEqual(requests, wantRequests) {
		t.Fatalf("requests = %v, want %v", requests, wantRequests)
	}
}

func TestUploadFileIncludesWebDAVResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			w.WriteHeader(http.StatusConflict)
			_, _ = io.WriteString(w, "<D:error>parent collection missing</D:error>")
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer server.Close()

	tempDir := t.TempDir()
	localPath := filepath.Join(tempDir, "backup.json")
	if err := os.WriteFile(localPath, []byte("backup"), 0600); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	client := NewClient(server.URL+"/dav/", "user", "password")
	err := client.UploadFile(localPath, "backup.json")
	if err == nil {
		t.Fatal("UploadFile returned nil for a 409 response")
	}
	if !strings.Contains(err.Error(), "status 409") || !strings.Contains(err.Error(), "parent collection missing") {
		t.Fatalf("UploadFile error = %q, want status and response body", err)
	}
}

func TestTestChecksCollectionWithoutUploading(t *testing.T) {
	var method string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if r.Method != "PROPFIND" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusMultiStatus)
		_, _ = io.WriteString(w, `<?xml version="1.0"?><D:multistatus xmlns:D="DAV:"/>`)
	}))
	defer server.Close()

	client := NewClient(server.URL+"/dav/", "user", "password")
	if err := client.Test(context.Background(), "/backup"); err != nil {
		t.Fatalf("Test returned error: %v", err)
	}
	if method != "PROPFIND" {
		t.Fatalf("Test used %q, want PROPFIND", method)
	}
}

func TestTestRejectsMissingCollection(t *testing.T) {
	var method string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if r.Method == "PROPFIND" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer server.Close()

	client := NewClient(server.URL+"/dav/", "user", "password")
	err := client.Test(context.Background(), "/new-backup")
	if err == nil {
		t.Fatal("Test returned nil for a missing collection")
	}
	if !strings.Contains(err.Error(), "status 404") {
		t.Fatalf("Test error = %q, want status 404", err)
	}
	if method != "PROPFIND" {
		t.Fatalf("Test used %q, want PROPFIND", method)
	}
}

func TestUploadFileRejectsUnsafePathBeforeDirectoryRequests(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	tempDir := t.TempDir()
	localPath := filepath.Join(tempDir, "backup.json")
	if err := os.WriteFile(localPath, []byte("backup"), 0600); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	client := NewClient(server.URL+"/dav/", "user", "password")
	if err := client.UploadFile(localPath, "nested/../backup.json"); err == nil {
		t.Fatal("UploadFile accepted an unsafe remote path")
	}
	if requests != 0 {
		t.Fatalf("UploadFile sent %d requests for an unsafe remote path, want 0", requests)
	}
}

func TestResponseErrorWithoutBody(t *testing.T) {
	resp := &http.Response{StatusCode: http.StatusConflict, Body: io.NopCloser(strings.NewReader(""))}
	err := responseError("upload", resp)
	if got, want := err.Error(), fmt.Sprintf("upload failed with status %d", http.StatusConflict); got != want {
		t.Fatalf("responseError() = %q, want %q", got, want)
	}
}
