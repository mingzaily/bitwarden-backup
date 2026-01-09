package provider

import "github.com/mingzaily/bitwarden-backup/internal/model"

// BackupContext 备份上下文，包含备份所需的所有信息
type BackupContext struct {
	SourceFile  string // 源文件路径
	TaskName    string // 任务名称
	Timestamp   string // 时间戳
	Destination model.BackupDestination
}

// DestinationProvider 备份目标提供者接口
type DestinationProvider interface {
	// Type 返回提供者类型标识
	Type() string

	// Backup 执行备份操作，返回最终存储路径
	Backup(ctx BackupContext) (string, error)
}

// RetentionProvider 支持备份保留策略的提供者接口
type RetentionProvider interface {
	// Cleanup 清理超出保留数量的旧备份
	// maxCount: 最大保留数量
	// 返回删除的文件数量和错误
	Cleanup(dest model.BackupDestination, maxCount int) (int, error)
}
