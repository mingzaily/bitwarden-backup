package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/bitwarden"
	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

// ServerProvider 服务器备份提供者
type ServerProvider struct{}

// NewServerProvider 创建服务器备份提供者
func NewServerProvider() *ServerProvider {
	return &ServerProvider{}
}

// Type 返回提供者类型
func (p *ServerProvider) Type() string {
	return "server"
}

// Backup 执行服务器备份（导入到目标服务器），返回目标服务器信息
func (p *ServerProvider) Backup(ctx BackupContext) (string, error) {
	dest := ctx.Destination

	if dest.TargetServerID == nil {
		ctx.AddLog("server", "服务器导入失败: target server id is nil")
		return "", fmt.Errorf("target server id is nil")
	}

	var targetServer model.ServerConfig
	if err := database.DB.First(&targetServer, *dest.TargetServerID).Error; err != nil {
		ctx.AddLog("server", "服务器导入失败: "+err.Error())
		return "", fmt.Errorf("failed to get target server: %w", err)
	}
	if !targetServer.Enabled {
		ctx.AddLog("server", fmt.Sprintf("目标服务器已停用: %s", targetServer.Name))
		return "", fmt.Errorf("target server is disabled: %s", targetServer.Name)
	}

	bwCtx := ctx.Context
	if bwCtx == nil {
		bwCtx = context.Background()
	}
	bwCtx, cancel := context.WithTimeout(bwCtx, 5*time.Minute)
	defer cancel()
	ctx.AddLog("server", fmt.Sprintf("开始导入服务器: %s", targetServer.Name))
	client := bitwarden.NewClientWithLogSink("server", ctx.Log)
	if err := client.WithProcessLock(bwCtx, func(lockedCtx context.Context) (err error) {
		// Clear any previous CLI session before switching the global server and
		// always clean up again on every early-return path.
		_ = client.Logout(lockedCtx)
		defer func() {
			cleanupCtx, cleanupCancel := context.WithTimeout(context.WithoutCancel(lockedCtx), 30*time.Second)
			defer cleanupCancel()
			if logoutErr := client.Logout(cleanupCtx); logoutErr != nil && err == nil {
				err = fmt.Errorf("failed to logout from target: %w", logoutErr)
			}
		}()

		if err := client.ConfigServer(lockedCtx, targetServer.ServerURL); err != nil {
			return fmt.Errorf("failed to config target server: %w", err)
		}
		if err := client.Login(lockedCtx, targetServer.ClientID, targetServer.ClientSecret); err != nil {
			return fmt.Errorf("failed to login to target: %w", err)
		}
		if err := client.Sync(lockedCtx); err != nil {
			return fmt.Errorf("failed to sync target: %w", err)
		}
		if err := client.Unlock(lockedCtx, targetServer.MasterPassword); err != nil {
			return fmt.Errorf("failed to unlock target: %w", err)
		}
		if err := client.Import(lockedCtx, ctx.SourceFile, "json"); err != nil {
			return fmt.Errorf("failed to import: %w", err)
		}
		return nil
	}); err != nil {
		ctx.AddLog("server", "服务器导入失败: "+err.Error())
		return "", err
	}

	ctx.AddLog("server", fmt.Sprintf("服务器导入完成: %s", targetServer.Name))
	// 返回目标服务器信息
	return fmt.Sprintf("server://%s", targetServer.Name), nil
}
