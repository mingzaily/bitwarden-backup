# Security Documentation

## Encryption Implementation

### Overview

This application implements **AES-256-GCM encryption** for all sensitive credentials stored in the SQLite3 database.

### Encrypted Fields

#### ServerConfig
- `ClientID` - Bitwarden API client identifier
- `ClientSecret` - Bitwarden API client secret
- `MasterPassword` - Bitwarden master password

#### BackupDestination
- `WebDAVPassword` - WebDAV server password

### Encryption Details

**Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Size**: 256 bits (32 bytes)
- **Authentication**: Built-in AEAD (Authenticated Encryption with Associated Data)
- **Nonce**: Random 12-byte nonce per encryption operation
- **Encoding**: Base64 for storage

### Key Management

#### Master Key Configuration

Set the encryption master key via environment variable:

```bash
export BITWARDEN_BACKUP_MASTER_KEY="your-secure-random-key-here"
```

**Important**:
- Use a strong, randomly generated key (minimum 32 characters)
- Store this key securely (password manager, secrets vault, etc.)
- **Never commit this key to version control**
- Losing this key means **permanent data loss**

#### Key Derivation

The master key is processed through PBKDF2:
- **Algorithm**: PBKDF2-HMAC-SHA256
- **Iterations**: 100,000
- **Salt**: Fixed application salt (consider using random salt in production)
- **Output**: 256-bit encryption key

### Automatic Encryption/Decryption

Encryption and decryption happen automatically via GORM hooks:

- **BeforeSave**: Encrypts sensitive fields before writing to database
- **AfterFind**: Decrypts sensitive fields after reading from database

No manual encryption/decryption needed in application code.

### Migration for Existing Data

If you have existing unencrypted data in your database, run the migration:

```go
import "github.com/mingzaily/bitwarden-backup/internal/database"

// Run once after enabling encryption
err := database.MigrateEncryptExistingData()
if err != nil {
    log.Fatal(err)
}
```

The migration:
- Detects unencrypted data automatically
- Encrypts all sensitive fields
- Skips already encrypted data
- Logs progress for each record

### API Response Masking

Sensitive data is masked in API responses for security:

**Original**: `client_secret_abc123xyz789`
**Masked**: `clie****z789`

This prevents credential leakage in:
- API responses
- Log files
- Browser developer tools
- Network traffic inspection

### Security Best Practices

1. **Master Key Security**
   - Generate using: `openssl rand -base64 32`
   - Store in secure secrets management system
   - Rotate periodically (requires re-encryption)
   - Use different keys for dev/staging/production

2. **Database File Protection**
   - Set restrictive file permissions: `chmod 600 bitwarden.db`
   - Store in secure location
   - Exclude from backups or encrypt backups
   - Never commit to version control

3. **Environment Variables**
   - Use `.env` files (add to `.gitignore`)
   - Or use system environment variables
   - Or use secrets management (HashiCorp Vault, AWS Secrets Manager)

4. **Backup Security**
   - Encrypted database is still sensitive (metadata visible)
   - Protect backup files with additional encryption
   - Secure backup storage locations
   - Regular backup testing

### Threat Model

**Protected Against**:
- ✅ Database file theft
- ✅ SQL injection credential exposure
- ✅ Memory dumps (credentials encrypted at rest)
- ✅ Accidental credential logging
- ✅ Unauthorized API access to credentials

**Not Protected Against**:
- ❌ Application memory inspection while running (decrypted in memory)
- ❌ Master key compromise
- ❌ Server-side code execution
- ❌ Physical access to running system

### Compliance

This implementation helps meet requirements for:
- **GDPR**: Encryption of personal data
- **PCI DSS**: Protection of authentication credentials
- **SOC 2**: Data encryption controls
- **HIPAA**: Encryption of sensitive information

### Troubleshooting

**Error: "encryption key not set"**
- Set `BITWARDEN_BACKUP_MASTER_KEY` environment variable
- Restart the application

**Error: "invalid ciphertext"**
- Master key changed or corrupted
- Database corruption
- Restore from backup with correct key

**Migration Issues**
- Check database permissions
- Verify master key is set
- Review migration logs for specific errors

### Technical Details

**Encryption Flow**:
```
Plaintext → AES-256-GCM → Ciphertext → Base64 → Database
Database → Base64 Decode → AES-256-GCM → Plaintext
```

**Storage Overhead**:
- Nonce: 12 bytes
- Auth tag: 16 bytes
- Base64 encoding: ~33% increase
- Total: ~50-60% larger than plaintext

**Performance**:
- Encryption: ~0.1ms per field
- Decryption: ~0.1ms per field
- Negligible impact on API response time
