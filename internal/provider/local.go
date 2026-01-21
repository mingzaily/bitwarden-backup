package provider

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// LocalProvider 本地存储提供者
type LocalProvider struct{}

// NewLocalProvider 创建本地存储提供者
func NewLocalProvider() *LocalProvider {
	return &LocalProvider{}
}

// Type 返回提供者类型
func (p *LocalProvider) Type() string {
	return "local"
}

// Backup 执行本地备份，返回最终存储路径
func (p *LocalProvider) Backup(ctx BackupContext) (string, error) {
	dest := ctx.Destination

	if err := os.MkdirAll(dest.LocalPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create local directory: %w", err)
	}

	targetFile := filepath.Join(dest.LocalPath, fmt.Sprintf("backup_%s_%s.json", ctx.TaskName, ctx.Timestamp))

	source, err := os.Open(ctx.SourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	target, err := os.Create(targetFile)
	if err != nil {
		return "", fmt.Errorf("failed to create target file: %w", err)
	}
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return targetFile, nil
}

// Cleanup 清理超出保留数量的旧备份
func (p *LocalProvider) Cleanup(dest model.BackupDestination, maxCount int) (int, error) {
	if maxCount <= 0 {
		return 0, nil
	}

	if dest.LocalPath == "" {
		return 0, fmt.Errorf("local path is empty")
	}

	entries, err := os.ReadDir(dest.LocalPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read directory: %w", err)
	}

	// 筛选备份文件并获取文件信息
	type backupFile struct {
		name    string
		modTime int64
	}
	var backups []backupFile

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, "backup_") || !strings.HasSuffix(name, ".json") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		backups = append(backups, backupFile{name: name, modTime: info.ModTime().UnixNano()})
	}

	if len(backups) <= maxCount {
		return 0, nil
	}

	// 按修改时间降序排序（最新在前）
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].modTime > backups[j].modTime
	})

	// 删除超出数量的旧文件
	deleted := 0
	for i := maxCount; i < len(backups); i++ {
		filePath := filepath.Join(dest.LocalPath, backups[i].name)
		if err := os.Remove(filePath); err != nil {
			continue
		}
		deleted++
	}

	return deleted, nil
}
