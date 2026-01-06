package database

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/mingzaily/bitwarden-backup/internal/crypto"
	"gorm.io/gorm"
)

// MigrateEncryptExistingData 迁移现有数据，加密敏感字段
// 此函数应该在首次启用加密时运行一次
func MigrateEncryptExistingData() error {
	log.Println("Starting migration: encrypting existing sensitive data...")

	// 迁移 ServerConfig 数据
	if err := migrateServerConfigs(); err != nil {
		return fmt.Errorf("failed to migrate server configs: %w", err)
	}

	// 迁移 BackupDestination 数据
	if err := migrateBackupDestinations(); err != nil {
		return fmt.Errorf("failed to migrate backup destinations: %w", err)
	}

	log.Println("Migration completed successfully")
	return nil
}

// migrateServerConfigs 迁移服务器配置数据
func migrateServerConfigs() error {
	var servers []ServerConfig

	// 直接从数据库读取，跳过 AfterFind 钩子
	if err := DB.Session(&gorm.Session{SkipHooks: true}).Find(&servers).Error; err != nil {
		return err
	}

	log.Printf("Found %d server configs to migrate", len(servers))

	for i := range servers {
		server := &servers[i]
		needsUpdate := false

		// 检查是否已加密（加密后的数据是 base64 编码，长度会更长）
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
			// 跳过 BeforeSave 钩子，直接保存已加密的数据
			if err := DB.Session(&gorm.Session{SkipHooks: true}).Save(server).Error; err != nil {
				return fmt.Errorf("failed to save server %d: %w", server.ID, err)
			}
			log.Printf("Migrated server config ID: %d", server.ID)
		}
	}

	return nil
}

// migrateBackupDestinations 迁移备份目标数据
func migrateBackupDestinations() error {
	var destinations []BackupDestination

	// 直接从数据库读取，跳过 AfterFind 钩子
	if err := DB.Session(&gorm.Session{SkipHooks: true}).Find(&destinations).Error; err != nil {
		return err
	}

	log.Printf("Found %d backup destinations to migrate", len(destinations))

	for i := range destinations {
		dest := &destinations[i]
		needsUpdate := false

		// 检查 WebDAV 密码是否需要加密
		if dest.WebDAVPassword != "" && !isEncrypted(dest.WebDAVPassword) {
			encrypted, err := crypto.Encrypt(dest.WebDAVPassword)
			if err != nil {
				return fmt.Errorf("failed to encrypt WebDAVPassword for destination %d: %w", dest.ID, err)
			}
			dest.WebDAVPassword = encrypted
			needsUpdate = true
		}

		if needsUpdate {
			// 跳过 BeforeSave 钩子，直接保存已加密的数据
			if err := DB.Session(&gorm.Session{SkipHooks: true}).Save(dest).Error; err != nil {
				return fmt.Errorf("failed to save destination %d: %w", dest.ID, err)
			}
			log.Printf("Migrated backup destination ID: %d", dest.ID)
		}
	}

	return nil
}

// isEncrypted 检查字符串是否已经加密（简单检查是否为有效的 base64）
func isEncrypted(s string) bool {
	// 加密后的数据是 base64 编码，尝试解码
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil && len(s) > 50 // 加密后的数据通常较长
}
