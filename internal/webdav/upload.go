package webdav

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/safety"
)

// httpClient 带超时的 HTTP 客户端
var httpClient = &http.Client{
	Timeout: 60 * time.Second,
}

const (
	maxPropfindResponse = 4 << 20
	maxErrorResponse    = 64 << 10
)

// ResponseError preserves the HTTP status returned by a WebDAV server so
// callers can distinguish an unreachable collection from a collection that
// simply has not been created yet.
type ResponseError struct {
	Action     string
	StatusCode int
	Detail     string
}

func (e *ResponseError) Error() string {
	message := fmt.Sprintf("%s failed with status %d", e.Action, e.StatusCode)
	if e.Detail != "" {
		message += ": " + e.Detail
	}
	return message
}

// responseError converts a non-success WebDAV response into a bounded error.
// Providers include this error in task execution logs, so keep the response
// body useful for diagnosis without allowing a remote server to flood logs.
func responseError(action string, resp *http.Response) error {
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxErrorResponse+1))
	detail := strings.Join(strings.Fields(string(body)), " ")
	if len(detail) > maxErrorResponse {
		detail = detail[:maxErrorResponse] + "…"
	}
	if readErr != nil {
		detail = fmt.Sprintf("read response: %v", readErr)
	}
	return &ResponseError{Action: action, StatusCode: resp.StatusCode, Detail: detail}
}

// ensureDirectory creates each missing collection in a remote path. WebDAV
// servers are not required to create intermediate collections as a side
// effect of PUT, and many providers return 409 when the parent is missing.
func (c *Client) ensureDirectory(ctx context.Context, remotePath string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	cleanPath := strings.Trim(remotePath, "/")
	if cleanPath == "" {
		return nil
	}

	current := ""
	for _, segment := range strings.Split(cleanPath, "/") {
		current = path.Join(current, segment)
		if err := c.ensureCollection(ctx, current); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ensureCollection(ctx context.Context, remotePath string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	collectionPath := strings.TrimSuffix(remotePath, "/") + "/"
	fullURL, err := c.requestURL(collectionPath)
	if err != nil {
		return err
	}

	probe, err := http.NewRequest("PROPFIND", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create directory check request: %w", err)
	}
	probe = probe.WithContext(ctx)
	probe.SetBasicAuth(c.username, c.password)
	probe.Header.Set("Depth", "0")

	resp, err := httpClient.Do(probe)
	if err != nil {
		return fmt.Errorf("failed to check WebDAV directory: %w", err)
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMultiStatus {
		resp.Body.Close()
		return nil
	}
	if resp.StatusCode != http.StatusNotFound {
		defer resp.Body.Close()
		return responseError("WebDAV directory check", resp)
	}
	resp.Body.Close()

	create, err := http.NewRequest("MKCOL", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create directory request: %w", err)
	}
	create = create.WithContext(ctx)
	create.SetBasicAuth(c.username, c.password)

	resp, err = httpClient.Do(create)
	if err != nil {
		return fmt.Errorf("failed to create WebDAV directory: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusMethodNotAllowed {
		// 405 commonly means another worker created the collection between the
		// PROPFIND and MKCOL requests. The following PUT will verify usability.
		return nil
	}
	return responseError("WebDAV directory creation", resp)
}

func (c *Client) requestURL(remotePath string) (string, error) {
	base, err := url.Parse(c.baseURL)
	if err != nil || base.Host == "" || (base.Scheme != "http" && base.Scheme != "https") {
		return "", fmt.Errorf("invalid WebDAV URL")
	}
	if err := safety.ValidateURL(c.baseURL, "webdav_url", true); err != nil {
		return "", err
	}
	if strings.ContainsAny(remotePath, "\r\n\\") {
		return "", fmt.Errorf("invalid WebDAV path")
	}
	for _, segment := range strings.Split(strings.Trim(remotePath, "/"), "/") {
		if segment == ".." {
			return "", fmt.Errorf("WebDAV path must not contain parent segments")
		}
	}
	joined := path.Join(base.Path, "/"+strings.TrimPrefix(remotePath, "/"))
	if strings.HasSuffix(remotePath, "/") && !strings.HasSuffix(joined, "/") {
		joined += "/"
	}
	base.Path = joined
	base.RawPath = ""
	return base.String(), nil
}

// UploadFile 上传文件到 WebDAV
func (c *Client) UploadFile(localPath, remotePath string) error {
	return c.UploadFileContext(context.Background(), localPath, remotePath)
}

// UploadFileContext 上传文件到 WebDAV，并将调用方的取消信号传递给目录
// 检查、目录创建和 PUT 请求。
func (c *Client) UploadFileContext(ctx context.Context, localPath, remotePath string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	// Validate the complete remote path before creating any parent collection.
	// This prevents an unsafe path from causing side effects during preparation.
	fullURL, err := c.requestURL(remotePath)
	if err != nil {
		return err
	}

	// 读取本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	parentPath := path.Dir(remotePath)
	if parentPath != "." && parentPath != "/" {
		if err := c.ensureDirectory(ctx, parentPath); err != nil {
			return fmt.Errorf("failed to prepare WebDAV directory %q: %w", parentPath, err)
		}
	}

	// 创建 PUT 请求
	req, err := http.NewRequest("PUT", fullURL, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)

	// 设置认证
	req.SetBasicAuth(c.username, c.password)
	req.ContentLength = fileInfo.Size()
	req.Header.Set("Content-Type", "application/octet-stream")

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return responseError("WebDAV upload", resp)
	}

	return nil
}

// FileInfo WebDAV 文件信息
type FileInfo struct {
	Name    string
	ModTime time.Time
	IsDir   bool
}

// multistatus PROPFIND 响应结构
type multistatus struct {
	Responses []response `xml:"response"`
}

type response struct {
	Href     string   `xml:"href"`
	Propstat propstat `xml:"propstat"`
}

type propstat struct {
	Prop prop `xml:"prop"`
}

type prop struct {
	DisplayName     string `xml:"displayname"`
	GetLastModified string `xml:"getlastmodified"`
	ResourceType    struct {
		Collection *struct{} `xml:"collection"`
	} `xml:"resourcetype"`
}

// ListFiles 列举目录下的文件
func (c *Client) ListFiles(remotePath string) ([]FileInfo, error) {
	return c.ListFilesContext(context.Background(), remotePath)
}

// ListFilesContext 列举目录下的文件，并允许调用方取消连接测试。
func (c *Client) ListFilesContext(ctx context.Context, remotePath string) ([]FileInfo, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	fullURL, err := c.requestURL(remotePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PROPFIND", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Depth", "1")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMultiStatus && resp.StatusCode != http.StatusOK {
		return nil, responseError("WebDAV list", resp)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxPropfindResponse+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if len(body) > maxPropfindResponse {
		return nil, fmt.Errorf("WebDAV response is too large")
	}

	var ms multistatus
	if err := xml.Unmarshal(body, &ms); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var files []FileInfo
	requestPath := strings.TrimPrefix(remotePath, "/")
	for _, r := range ms.Responses {
		// 跳过目录本身（基于 Href 判断）
		href := strings.TrimSuffix(r.Href, "/")
		if strings.HasSuffix(href, requestPath) || href == "" {
			continue
		}
		isDir := r.Propstat.Prop.ResourceType.Collection != nil
		modTime := parseWebDAVTime(r.Propstat.Prop.GetLastModified)
		name := r.Propstat.Prop.DisplayName
		if name == "" {
			parts := strings.Split(href, "/")
			name = parts[len(parts)-1]
		}
		files = append(files, FileInfo{Name: name, ModTime: modTime, IsDir: isDir})
	}

	return files, nil
}

// Test checks the configured WebDAV collection without uploading or deleting
// any file. A non-2xx response includes its bounded response body in the
// returned error for the task execution log/UI.
func (c *Client) Test(ctx context.Context, remotePath string) error {
	_, err := c.ListFilesContext(ctx, remotePath)
	if err == nil {
		return nil
	}
	return err
}

// Delete 删除远程文件
func (c *Client) Delete(remotePath string) error {
	return c.DeleteContext(context.Background(), remotePath)
}

// DeleteContext 删除远程文件，并将调用方的取消信号传递给请求。
func (c *Client) DeleteContext(ctx context.Context, remotePath string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	fullURL, err := c.requestURL(remotePath)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req = req.WithContext(ctx)

	req.SetBasicAuth(c.username, c.password)

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return responseError("WebDAV delete", resp)
	}

	return nil
}

// parseWebDAVTime 解析 WebDAV 时间格式
func parseWebDAVTime(s string) time.Time {
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		"Mon, 02 Jan 2006 15:04:05 GMT",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
