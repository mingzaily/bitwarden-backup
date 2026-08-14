package provider

import (
	"context"

	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// BackupContext 备份上下文，包含备份所需的所有信息
type BackupContext struct {
	Context          context.Context
	SourceFile       string // 源文件路径
	TaskName         string // 任务名称
	Timestamp        string // YYYYMMDDHHmmss
	FilenameTemplate string // 备份文件名模板
	Destination      model.BackupDestination
	Log              func(source, message string)
}

// AddLog writes a provider-scoped execution log when the scheduler supplied
// a sink. Providers remain usable in isolation when no sink is configured.
func (c BackupContext) AddLog(source, message string) {
	if c.Log != nil {
		c.Log(source, message)
	}
}

// DestinationProvider 备份目标提供者接口
type DestinationProvider interface {
	// Type 返回提供者类型标识
	Type() string

	// Backup 执行备份操作，返回最终存储路径
	Backup(ctx BackupContext) (string, error)
}

// ConnectionTester is implemented by providers that can validate their
// connection without creating a backup file.
type ConnectionTester interface {
	Test(ctx context.Context, destination model.BackupDestination) error
}

// RetentionProvider 支持备份保留策略的提供者接口
type RetentionProvider interface {
	// Cleanup 清理当前任务文件名模板下超出保留数量的旧备份
	// maxCount: 最大保留数量
	// 返回删除的文件数量和错误
	Cleanup(ctx BackupContext, maxCount int) (int, error)
}
