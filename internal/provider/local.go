package provider

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
	ctx.AddLog("local", fmt.Sprintf("开始本地 CP: %s", dest.LocalPath))
	fail := func(err error) (string, error) {
		ctx.AddLog("local", "本地 CP 失败: "+err.Error())
		return "", err
	}

	if dest.LocalPath == "" {
		return fail(fmt.Errorf("local path is empty"))
	}
	if err := os.MkdirAll(dest.LocalPath, 0700); err != nil {
		return fail(fmt.Errorf("failed to create local directory: %w", err))
	}
	if err := os.Chmod(dest.LocalPath, 0700); err != nil {
		return fail(fmt.Errorf("failed to secure local directory: %w", err))
	}

	filename := renderBackupFilename(ctx)
	targetFile := filepath.Join(dest.LocalPath, filename)

	source, err := os.Open(ctx.SourceFile)
	if err != nil {
		return fail(fmt.Errorf("failed to open source file: %w", err))
	}
	defer source.Close()

	target, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fail(fmt.Errorf("failed to create target file: %w", err))
	}
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return fail(fmt.Errorf("failed to copy file: %w", err))
	}

	ctx.AddLog("local", fmt.Sprintf("本地 CP 完成: %s", targetFile))
	return targetFile, nil
}

// Cleanup 清理超出保留数量的旧备份
func (p *LocalProvider) Cleanup(ctx BackupContext, maxCount int) (int, error) {
	if maxCount <= 0 {
		return 0, nil
	}

	dest := ctx.Destination
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
	var inspectErrors []error

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !matchesBackupFilename(name, ctx) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			inspectErrors = append(inspectErrors, fmt.Errorf("%s: %w", name, err))
			continue
		}
		backups = append(backups, backupFile{name: name, modTime: info.ModTime().UnixNano()})
	}
	if len(inspectErrors) > 0 {
		return 0, fmt.Errorf("failed to inspect local backup files: %w", errors.Join(inspectErrors...))
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
	var deleteErrors []error
	for i := maxCount; i < len(backups); i++ {
		filePath := filepath.Join(dest.LocalPath, backups[i].name)
		if err := os.Remove(filePath); err != nil {
			deleteErrors = append(deleteErrors, fmt.Errorf("%s: %w", backups[i].name, err))
			continue
		}
		deleted++
	}
	if len(deleteErrors) > 0 {
		return deleted, fmt.Errorf("failed to remove old local backups: %w", errors.Join(deleteErrors...))
	}

	return deleted, nil
}
