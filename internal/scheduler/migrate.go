package scheduler

import (
	"fmt"

	"github.com/mingzaily/bitwarden-backup/internal/bitwarden"
	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func (s *Scheduler) performMigration(client *bitwarden.Client, targetServerID uint, backupFile string) error {
	var targetServer model.ServerConfig
	if err := database.DB.First(&targetServer, targetServerID).Error; err != nil {
		return fmt.Errorf("failed to get target server: %w", err)
	}

	if err := client.Logout(); err != nil {
		return fmt.Errorf("failed to logout from source: %w", err)
	}

	if err := client.ConfigServer(targetServer.ServerURL); err != nil {
		return fmt.Errorf("failed to config target server: %w", err)
	}

	if err := client.Login(targetServer.ClientID, targetServer.ClientSecret); err != nil {
		return fmt.Errorf("failed to login to target: %w", err)
	}

	if err := client.Unlock(targetServer.MasterPassword); err != nil {
		return fmt.Errorf("failed to unlock target: %w", err)
	}

	if err := client.Import(backupFile, "json"); err != nil {
		return fmt.Errorf("failed to import: %w", err)
	}

	return nil
}
