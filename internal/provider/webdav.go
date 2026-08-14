package provider

import (
	"context"
	"errors"
	"fmt"
	"path"
	"sort"

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
	remoteFile := path.Join(dest.WebDAVPath, renderBackupFilename(ctx))
	ctx.AddLog("webdav", fmt.Sprintf("开始 WebDAV 上传: %s", remoteFile))

	if err := client.UploadFileContext(ctx.Context, ctx.SourceFile, remoteFile); err != nil {
		ctx.AddLog("webdav", "WebDAV 上传失败: "+err.Error())
		return "", fmt.Errorf("failed to upload to webdav: %w", err)
	}

	ctx.AddLog("webdav", fmt.Sprintf("WebDAV 上传完成: %s", remoteFile))
	// 返回完整的 WebDAV 路径
	return dest.WebDAVURL + remoteFile, nil
}

// Test verifies that the configured collection can be queried without
// uploading a backup file or creating a test task.
func (p *WebDAVProvider) Test(ctx context.Context, dest model.BackupDestination) error {
	client := webdav.NewClient(dest.WebDAVURL, dest.WebDAVUsername, dest.WebDAVPassword)
	if err := client.Test(ctx, dest.WebDAVPath); err != nil {
		return fmt.Errorf("WebDAV connection test failed: %w", err)
	}
	return nil
}

// Cleanup 清理超出保留数量的旧备份
func (p *WebDAVProvider) Cleanup(ctx BackupContext, maxCount int) (int, error) {
	if maxCount <= 0 {
		return 0, nil
	}

	dest := ctx.Destination
	client := webdav.NewClient(dest.WebDAVURL, dest.WebDAVUsername, dest.WebDAVPassword)
	files, err := client.ListFilesContext(ctx.Context, dest.WebDAVPath)
	if err != nil {
		return 0, fmt.Errorf("failed to list files: %w", err)
	}

	// 筛选备份文件
	var backups []webdav.FileInfo
	for _, f := range files {
		if f.IsDir {
			continue
		}
		if !matchesBackupFilename(f.Name, ctx) {
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
	var deleteErrors []error
	for i := maxCount; i < len(backups); i++ {
		remotePath := path.Join(dest.WebDAVPath, backups[i].Name)
		if err := client.DeleteContext(ctx.Context, remotePath); err != nil {
			deleteErrors = append(deleteErrors, fmt.Errorf("%s: %w", backups[i].Name, err))
			continue
		}
		deleted++
	}
	if len(deleteErrors) > 0 {
		return deleted, fmt.Errorf("failed to remove old WebDAV backups: %w", errors.Join(deleteErrors...))
	}

	return deleted, nil
}
