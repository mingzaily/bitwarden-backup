package scheduler

import (
	"fmt"

	"github.com/mingzaily/bitwarden-backup/internal/bitwarden"
	"github.com/mingzaily/bitwarden-backup/internal/database"
	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func (s *Scheduler) backupToServer(dest model.BackupDestination, sourceFile string) error {
	if dest.TargetServerID == nil {
		return fmt.Errorf("target server id is nil")
	}

	var targetServer model.ServerConfig
	if err := database.DB.First(&targetServer, *dest.TargetServerID).Error; err != nil {
		return fmt.Errorf("failed to get target server: %w", err)
	}

	client := bitwarden.NewClient()
	if err := client.ConfigServer(targetServer.ServerURL); err != nil {
		return fmt.Errorf("failed to config target server: %w", err)
	}

	if err := client.Login(targetServer.ClientID, targetServer.ClientSecret); err != nil {
		return fmt.Errorf("failed to login to target: %w", err)
	}

	if err := client.Unlock(targetServer.MasterPassword); err != nil {
		return fmt.Errorf("failed to unlock target: %w", err)
	}

	if err := client.Import(sourceFile, "json"); err != nil {
		client.Logout()
		return fmt.Errorf("failed to import: %w", err)
	}

	client.Logout()
	return nil
}
