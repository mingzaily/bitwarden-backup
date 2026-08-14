package webdav

import (
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

const maxPropfindResponse = 4 << 20

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
	// 读取本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 构建完整的远程路径
	fullURL, err := c.requestURL(remotePath)
	if err != nil {
		return err
	}

	// 创建 PUT 请求
	req, err := http.NewRequest("PUT", fullURL, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置认证
	req.SetBasicAuth(c.username, c.password)

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("upload failed with status: %d", resp.StatusCode)
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
	fullURL, err := c.requestURL(remotePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PROPFIND", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Depth", "1")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMultiStatus && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list failed with status: %d", resp.StatusCode)
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

// Delete 删除远程文件
func (c *Client) Delete(remotePath string) error {
	fullURL, err := c.requestURL(remotePath)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete failed with status: %d", resp.StatusCode)
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
