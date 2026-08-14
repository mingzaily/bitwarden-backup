package model

import "github.com/mingzaily/bitwarden-backup/internal/crypto"

// decryptSensitive supports both the current prefixed ciphertext and the
// legacy unprefixed AES-GCM format. Plaintext values are returned unchanged so
// an interrupted migration remains readable and can be upgraded on startup.
func decryptSensitive(value string) (string, error) {
	if crypto.IsEncrypted(value) {
		return crypto.Decrypt(value)
	}
	if decrypted, err := crypto.Decrypt(value); err == nil {
		return decrypted, nil
	}
	return value, nil
}
