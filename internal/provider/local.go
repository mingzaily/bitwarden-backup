package provider

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// Backup 执行本地备份
func (p *LocalProvider) Backup(ctx BackupContext) error {
	dest := ctx.Destination

	if err := os.MkdirAll(dest.LocalPath, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	targetFile := filepath.Join(dest.LocalPath, fmt.Sprintf("backup_%s_%s.json", ctx.TaskName, ctx.Timestamp))

	source, err := os.Open(ctx.SourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	target, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
