package scheduler

import (
	"context"
	"fmt"

	"github.com/mingzaily/bitwarden-backup/internal/bitwarden"
	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func (s *Scheduler) performMigration(client *bitwarden.Client, targetServerID uint, backupFile string) error {
	ctx := context.Background()
	var targetServer model.ServerConfig
	if err := database.DB.First(&targetServer, targetServerID).Error; err != nil {
		return fmt.Errorf("failed to get target server: %w", err)
	}

	if err := client.Logout(ctx); err != nil {
		return fmt.Errorf("failed to logout from source: %w", err)
	}

	if err := client.ConfigServer(ctx, targetServer.ServerURL); err != nil {
		return fmt.Errorf("failed to config target server: %w", err)
	}

	if err := client.Login(ctx, targetServer.ClientID, targetServer.ClientSecret); err != nil {
		return fmt.Errorf("failed to login to target: %w", err)
	}

	if err := client.Sync(ctx); err != nil {
		return fmt.Errorf("failed to sync target: %w", err)
	}

	if err := client.Unlock(ctx, targetServer.MasterPassword); err != nil {
		return fmt.Errorf("failed to unlock target: %w", err)
	}

	if err := client.Import(ctx, backupFile, "json"); err != nil {
		return fmt.Errorf("failed to import: %w", err)
	}

	return nil
}
