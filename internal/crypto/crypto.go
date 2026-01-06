package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

var (
	// ErrInvalidCiphertext 无效的密文
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	// ErrEncryptionKeyNotSet 加密密钥未设置
	ErrEncryptionKeyNotSet = errors.New("encryption key not set")
)

// encryptionKey 全局加密密钥
var encryptionKey []byte

// InitEncryption 初始化加密系统
func InitEncryption() error {
	// 从环境变量获取主密钥
	masterKey := os.Getenv("BITWARDEN_BACKUP_MASTER_KEY")
	if masterKey == "" {
		// 如果未设置，生成一个随机密钥（仅用于开发环境）
		masterKey = generateRandomKey()
		// 在生产环境中，应该要求用户设置此环境变量
		// return ErrEncryptionKeyNotSet
	}

	// 使用 PBKDF2 派生加密密钥
	salt := []byte("bitwarden-backup-salt-v1") // 在生产环境中应该使用随机盐并存储
	encryptionKey = pbkdf2.Key([]byte(masterKey), salt, 100000, 32, sha256.New)

	return nil
}

// generateRandomKey 生成随机密钥（仅用于开发）
func generateRandomKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(key)
}

// Encrypt 加密数据
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
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密数据
func Decrypt(ciphertext string) (string, error) {
	if len(encryptionKey) == 0 {
		return "", ErrEncryptionKeyNotSet
	}

	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
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
