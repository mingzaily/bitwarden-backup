package provider

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mingzaily/bitwarden-backup/internal/model"
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

// Cleanup 清理超出保留数量的旧备份
func (p *WebDAVProvider) Cleanup(dest model.BackupDestination, maxCount int) (int, error) {
	if maxCount <= 0 {
		return 0, nil
	}

	client := webdav.NewClient(dest.WebDAVURL, dest.WebDAVUsername, dest.WebDAVPassword)
	files, err := client.ListFiles(dest.WebDAVPath)
	if err != nil {
		return 0, fmt.Errorf("failed to list files: %w", err)
	}

	// 筛选备份文件
	var backups []webdav.FileInfo
	for _, f := range files {
		if f.IsDir {
			continue
		}
		if !strings.HasPrefix(f.Name, "backup_") || !strings.HasSuffix(f.Name, ".json") {
			continue
		}
		backups = append(backups, f)
	}

	if len(backups) <= maxCount {
		return 0, nil
	}

	// 按修改时间降序排序
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].ModTime.After(backups[j].ModTime)
	})

	// 删除超出数量的旧文件
	deleted := 0
	for i := maxCount; i < len(backups); i++ {
		remotePath := filepath.Join(dest.WebDAVPath, backups[i].Name)
		if err := client.Delete(remotePath); err != nil {
			continue
		}
		deleted++
	}

	return deleted, nil
}
