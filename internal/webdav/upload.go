package webdav

import (
	"fmt"
	"net/http"
	"os"
)

// UploadFile 上传文件到 WebDAV
func (c *Client) UploadFile(localPath, remotePath string) error {
	// 读取本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 构建完整的远程路径
	fullURL := c.baseURL + "/" + remotePath

	// 创建 PUT 请求
	req, err := http.NewRequest("PUT", fullURL, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置认证
	req.SetBasicAuth(c.username, c.password)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
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
