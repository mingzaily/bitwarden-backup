package database

import (
	"encoding/base64"
	"fmt"
	"github.com/mingzaily/bitwarden-backup/internal/logger"

	"github.com/mingzaily/bitwarden-backup/internal/crypto"
	"github.com/mingzaily/bitwarden-backup/internal/model"
	"gorm.io/gorm"
)

func MigrateEncryptExistingData() error {
	logger.Module(logger.ModuleDatabase).Info("Starting migration: encrypting existing sensitive data...")

	if err := migrateServerConfigs(); err != nil {
		return fmt.Errorf("failed to migrate server configs: %w", err)
	}

	if err := migrateBackupDestinations(); err != nil {
		return fmt.Errorf("failed to migrate backup destinations: %w", err)
	}

	logger.Module(logger.ModuleDatabase).Info("Migration completed successfully")
	return nil
}

func migrateServerConfigs() error {
	var servers []model.ServerConfig

	if err := DB.Session(&gorm.Session{SkipHooks: true}).Find(&servers).Error; err != nil {
		return err
	}

	logger.Module(logger.ModuleDatabase).Info("Found server configs to migrate", "count", len(servers))

	for i := range servers {
		server := &servers[i]
		needsUpdate := false

		if server.ClientID != "" && !isEncrypted(server.ClientID) {
			encrypted, err := crypto.Encrypt(server.ClientID)
			if err != nil {
				return fmt.Errorf("failed to encrypt ClientID for server %d: %w", server.ID, err)
			}
			server.ClientID = encrypted
			needsUpdate = true
		}

		if server.ClientSecret != "" && !isEncrypted(server.ClientSecret) {
			encrypted, err := crypto.Encrypt(server.ClientSecret)
			if err != nil {
				return fmt.Errorf("failed to encrypt ClientSecret for server %d: %w", server.ID, err)
			}
			server.ClientSecret = encrypted
			needsUpdate = true
		}

		if server.MasterPassword != "" && !isEncrypted(server.MasterPassword) {
			encrypted, err := crypto.Encrypt(server.MasterPassword)
			if err != nil {
				return fmt.Errorf("failed to encrypt MasterPassword for server %d: %w", server.ID, err)
			}
			server.MasterPassword = encrypted
			needsUpdate = true
		}

		if needsUpdate {
			if err := DB.Session(&gorm.Session{SkipHooks: true}).Save(server).Error; err != nil {
				return fmt.Errorf("failed to save server %d: %w", server.ID, err)
			}
			logger.Module(logger.ModuleDatabase).Info("Migrated server config", "id", server.ID)
		}
	}

	return nil
}

func migrateBackupDestinations() error {
	var destinations []model.BackupDestination

	if err := DB.Session(&gorm.Session{SkipHooks: true}).Find(&destinations).Error; err != nil {
		return err
	}

	logger.Module(logger.ModuleDatabase).Info("Found backup destinations to migrate", "count", len(destinations))

	for i := range destinations {
		dest := &destinations[i]
		needsUpdate := false

		if dest.WebDAVPassword != "" && !isEncrypted(dest.WebDAVPassword) {
			encrypted, err := crypto.Encrypt(dest.WebDAVPassword)
			if err != nil {
				return fmt.Errorf("failed to encrypt WebDAVPassword for destination %d: %w", dest.ID, err)
			}
			dest.WebDAVPassword = encrypted
			needsUpdate = true
		}

		if needsUpdate {
			if err := DB.Session(&gorm.Session{SkipHooks: true}).Save(dest).Error; err != nil {
				return fmt.Errorf("failed to save destination %d: %w", dest.ID, err)
			}
			logger.Module(logger.ModuleDatabase).Info("Migrated backup destination", "id", dest.ID)
		}
	}

	return nil
}

func isEncrypted(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil && len(s) > 50
}
