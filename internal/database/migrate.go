package database

import (
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

		if encrypted, changed, err := migrateSensitiveValue(server.ClientID); err != nil {
			return fmt.Errorf("failed to encrypt ClientID for server %d: %w", server.ID, err)
		} else if changed {
			server.ClientID = encrypted
			needsUpdate = true
		}

		if encrypted, changed, err := migrateSensitiveValue(server.ClientSecret); err != nil {
			return fmt.Errorf("failed to encrypt ClientSecret for server %d: %w", server.ID, err)
		} else if changed {
			server.ClientSecret = encrypted
			needsUpdate = true
		}

		if encrypted, changed, err := migrateSensitiveValue(server.MasterPassword); err != nil {
			return fmt.Errorf("failed to encrypt MasterPassword for server %d: %w", server.ID, err)
		} else if changed {
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

		fields := []struct {
			name  string
			value *string
		}{
			{"WebDAVPassword", &dest.WebDAVPassword},
			{"S3AccessKey", &dest.S3AccessKey},
			{"S3SecretKey", &dest.S3SecretKey},
			{"EncryptionPassword", &dest.EncryptionPassword},
		}
		for _, field := range fields {
			encrypted, changed, err := migrateSensitiveValue(*field.value)
			if err != nil {
				return fmt.Errorf("failed to encrypt %s for destination %d: %w", field.name, dest.ID, err)
			}
			if changed {
				*field.value = encrypted
				needsUpdate = true
			}
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

func migrateSensitiveValue(value string) (string, bool, error) {
	if value == "" || crypto.IsEncrypted(value) {
		return value, false, nil
	}
	// Decrypt succeeds for the legacy unprefixed ciphertext. If it fails, the
	// value is treated as legacy plaintext and encrypted as-is.
	if plaintext, err := crypto.Decrypt(value); err == nil {
		value = plaintext
	}
	encrypted, err := crypto.Encrypt(value)
	if err != nil {
		return "", false, err
	}
	return encrypted, true, nil
}
