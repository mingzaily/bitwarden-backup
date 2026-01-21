package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"github.com/mingzaily/bitwarden-backup/internal/logger"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/pbkdf2"
)

var (
	// ErrInvalidCiphertext 无效的密文
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	// ErrEncryptionKeyNotSet 加密密钥未设置
	ErrEncryptionKeyNotSet = errors.New("encryption key not set")
)

// EncryptedPrefix 加密数据前缀标记，用于识别已加密的字段
const EncryptedPrefix = "enc:v1:"

// encryptionKey 全局加密密钥
var encryptionKey []byte

// InitEncryption 初始化加密系统
func InitEncryption() error {
	// 1. 首先检查环境变量
	masterKey := os.Getenv("BITWARDEN_BACKUP_MASTER_KEY")
	if masterKey != "" {
		logger.Module(logger.ModuleEncryption).Info("Master key loaded from environment variable")
		return deriveEncryptionKey(masterKey)
	}

	// 2. 尝试从 data/.env 文件加载（Docker 持久化目录）
	dataEnvPath := "data/.env"
	masterKey, err := loadKeyFromEnvFile(dataEnvPath)
	if err == nil && masterKey != "" {
		logger.Module(logger.ModuleEncryption).Info("Master key loaded from file", "file", dataEnvPath)
		return deriveEncryptionKey(masterKey)
	}

	// 3. 尝试从 .env 文件加载（本地开发兼容）
	envPath := ".env"
	masterKey, err = loadKeyFromEnvFile(envPath)
	if err == nil && masterKey != "" {
		logger.Module(logger.ModuleEncryption).Info("Master key loaded from file", "file", envPath)
		return deriveEncryptionKey(masterKey)
	}

	// 4. 生成新密钥并保存到 data/.env（优先）或 .env
	logger.Module(logger.ModuleEncryption).Info("No master key found, generating new key")

	// 确保 data 目录存在
	if err := os.MkdirAll("data", 0755); err == nil {
		masterKey, err = generateAndSaveKey(dataEnvPath)
		if err == nil {
			logger.Module(logger.ModuleEncryption).Info("New master key generated and saved", "file", dataEnvPath)
			logger.Module(logger.ModuleEncryption).Info("Key is persisted in data/ directory")
			return deriveEncryptionKey(masterKey)
		}
	}

	// 回退到 .env
	masterKey, err = generateAndSaveKey(envPath)
	if err != nil {
		return fmt.Errorf("failed to generate and save master key: %w", err)
	}

	logger.Module(logger.ModuleEncryption).Info("New master key generated and saved", "file", envPath)
	logger.Module(logger.ModuleEncryption).Info("IMPORTANT: Backup this .env file!")

	return deriveEncryptionKey(masterKey)
}

// deriveEncryptionKey 从主密钥派生加密密钥
// 使用固定盐是设计决策：主密钥本身是随机 32 字节，每个部署实例不同
func deriveEncryptionKey(masterKey string) error {
	// 使用 PBKDF2 派生加密密钥
	salt := []byte("bitwarden-backup-salt-v1")
	encryptionKey = pbkdf2.Key([]byte(masterKey), salt, 100000, 32, sha256.New)
	return nil
}

// loadKeyFromEnvFile 从 .env 文件加载密钥
func loadKeyFromEnvFile(envPath string) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return "", err
	}

	// 加载 .env 文件
	envMap, err := godotenv.Read(envPath)
	if err != nil {
		return "", fmt.Errorf("failed to read .env file: %w", err)
	}

	// 获取密钥
	masterKey, exists := envMap["BITWARDEN_BACKUP_MASTER_KEY"]
	if !exists || masterKey == "" {
		return "", errors.New("BITWARDEN_BACKUP_MASTER_KEY not found in .env")
	}

	return masterKey, nil
}

// generateRandomKey 生成随机密钥，返回 error 而非 panic
func generateRandomKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// generateAndSaveKey 生成新密钥并保存到 .env 文件
func generateAndSaveKey(envPath string) (string, error) {
	// 生成新密钥
	masterKey, err := generateRandomKey()
	if err != nil {
		return "", err
	}

	// 准备 .env 内容
	envMap := make(map[string]string)

	// 如果文件已存在，先读取现有内容（保护其他配置项）
	if _, err := os.Stat(envPath); err == nil {
		existingEnv, err := godotenv.Read(envPath)
		if err != nil {
			// 如果无法读取现有文件，返回错误而不是覆盖
			return "", fmt.Errorf("existing .env file is corrupted or unreadable: %w", err)
		}
		envMap = existingEnv
	}

	// 设置新密钥
	envMap["BITWARDEN_BACKUP_MASTER_KEY"] = masterKey

	// 使用安全的方式写入 .env 文件（创建时就设置 0600 权限）
	if err := writeEnvFileSecurely(envPath, envMap); err != nil {
		return "", fmt.Errorf("failed to write .env file: %w", err)
	}

	return masterKey, nil
}

// writeEnvFileSecurely 以安全的方式写入 .env 文件（创建时就设置 0600 权限）
func writeEnvFileSecurely(envPath string, envMap map[string]string) error {
	// 创建临时文件，直接设置 0600 权限
	tmpFile, err := os.OpenFile(envPath+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// 写入内容
	for key, value := range envMap {
		if _, err := fmt.Fprintf(tmpFile, "%s=\"%s\"\n", key, value); err != nil {
			os.Remove(envPath + ".tmp")
			return fmt.Errorf("failed to write to temp file: %w", err)
		}
	}

	// 确保数据写入磁盘
	if err := tmpFile.Sync(); err != nil {
		os.Remove(envPath + ".tmp")
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	// 关闭临时文件
	tmpFile.Close()

	// 原子性地重命名临时文件为目标文件
	if err := os.Rename(envPath+".tmp", envPath); err != nil {
		os.Remove(envPath + ".tmp")
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// Encrypt 加密数据，返回带前缀的密文
func Encrypt(plaintext string) (string, error) {
	if len(encryptionKey) == 0 {
		return "", ErrEncryptionKeyNotSet
	}

	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	// 添加前缀标记，便于识别已加密数据
	return EncryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密数据，支持带前缀和不带前缀的密文（向后兼容）
func Decrypt(ciphertext string) (string, error) {
	if len(encryptionKey) == 0 {
		return "", ErrEncryptionKeyNotSet
	}

	if ciphertext == "" {
		return "", nil
	}

	// 去除前缀标记（如果存在）
	encodedData := ciphertext
	if strings.HasPrefix(ciphertext, EncryptedPrefix) {
		encodedData = strings.TrimPrefix(ciphertext, EncryptedPrefix)
	}

	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// IsEncrypted 检查字符串是否已加密（通过前缀标记识别）
func IsEncrypted(s string) bool {
	return strings.HasPrefix(s, EncryptedPrefix)
}
