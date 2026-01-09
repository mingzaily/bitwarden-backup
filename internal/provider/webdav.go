package provider

import (
	"fmt"
	"path/filepath"

	"github.com/mingzaily/bitwarden-backup/internal/webdav"
)

// WebDAVProvider WebDAV 存储提供者
type WebDAVProvider struct{}

// NewWebDAVProvider 创建 WebDAV 存储提供者
func NewWebDAVProvider() *WebDAVProvider {
	return &WebDAVProvider{}
}

// Type 返回提供者类型
func (p *WebDAVProvider) Type() string {
	return "webdav"
}

// Backup 执行 WebDAV 备份，返回最终存储路径
func (p *WebDAVProvider) Backup(ctx BackupContext) (string, error) {
	dest := ctx.Destination

	client := webdav.NewClient(dest.WebDAVURL, dest.WebDAVUsername, dest.WebDAVPassword)
	remoteFile := filepath.Join(dest.WebDAVPath, fmt.Sprintf("backup_%s_%s.json", ctx.TaskName, ctx.Timestamp))

	if err := client.UploadFile(ctx.SourceFile, remoteFile); err != nil {
		return "", fmt.Errorf("failed to upload to webdav: %w", err)
	}

	// 返回完整的 WebDAV 路径
	return dest.WebDAVURL + remoteFile, nil
}
