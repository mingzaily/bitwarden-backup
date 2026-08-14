package handler

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/safety"
	"github.com/robfig/cron/v3"
)

var s3BucketPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)

func validateServerRequest(req model.ServerRequest, requireSecrets bool) error {
	if err := safety.ValidateName(req.Name, "name", 100); err != nil {
		return err
	}
	if err := safety.ValidateURL(req.ServerURL, "server_url", true); err != nil {
		return err
	}
	if err := validateText(req.ClientID, "client_id", 500, false); err != nil {
		return err
	}
	if requireSecrets && strings.TrimSpace(req.ClientSecret) == "" {
		return fmt.Errorf("client_secret is required")
	}
	if requireSecrets && strings.TrimSpace(req.MasterPassword) == "" {
		return fmt.Errorf("master_password is required")
	}
	if req.ClientSecret != "" {
		if err := validateText(req.ClientSecret, "client_secret", 500, false); err != nil {
			return err
		}
	}
	if req.MasterPassword != "" {
		if err := validateText(req.MasterPassword, "master_password", 500, false); err != nil {
			return err
		}
	}
	return nil
}

func validateDestination(dest model.DestinationRequest) error {
	if err := safety.ValidateName(dest.Name, "name", 100); err != nil {
		return err
	}
	if dest.MaxBackupCount < 0 {
		return fmt.Errorf("max_backup_count must not be negative")
	}

	switch dest.Type {
	case "local":
		if !filepath.IsAbs(strings.TrimSpace(dest.LocalPath)) {
			return fmt.Errorf("local_path must be an absolute path")
		}
		if err := validateText(dest.LocalPath, "local_path", 255, false); err != nil {
			return err
		}
	case "webdav":
		if err := safety.ValidateURL(dest.WebDAVURL, "webdav_url", true); err != nil {
			return err
		}
		if err := safety.ValidateRemotePath(dest.WebDAVPath, "webdav_path"); err != nil {
			return err
		}
		if err := validateText(dest.WebDAVUsername, "webdav_username", 100, false); err != nil {
			return err
		}
		if err := validateText(dest.WebDAVPassword, "webdav_password", 500, false); err != nil {
			return err
		}
	case "s3":
		if dest.S3Endpoint != "" {
			if err := safety.ValidateURL(dest.S3Endpoint, "s3_endpoint", true); err != nil {
				return err
			}
		}
		if !s3BucketPattern.MatchString(strings.ToLower(strings.TrimSpace(dest.S3Bucket))) {
			return fmt.Errorf("s3_bucket is invalid")
		}
		if err := safety.ValidateRemotePath(dest.S3Path, "s3_path"); err != nil {
			return err
		}
		if err := validateText(dest.S3Region, "s3_region", 100, false); err != nil {
			return err
		}
		if err := validateText(dest.S3AccessKey, "s3_access_key", 500, false); err != nil {
			return err
		}
		if err := validateText(dest.S3SecretKey, "s3_secret_key", 500, false); err != nil {
			return err
		}
	case "server":
		if dest.TargetServerID == nil || *dest.TargetServerID == 0 {
			return fmt.Errorf("target_server_id is required")
		}
	default:
		return fmt.Errorf("unsupported destination type")
	}

	if dest.Encrypted && dest.Type == "server" {
		return fmt.Errorf("server destinations do not support encrypted exports")
	}
	if dest.EncryptionPassword != "" {
		if err := validateText(dest.EncryptionPassword, "encryption_password", 500, false); err != nil {
			return err
		}
	}
	return nil
}

func validateText(value, field string, max int, required bool) error {
	if required && strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	if len([]rune(value)) > max {
		return fmt.Errorf("%s is too long", field)
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return fmt.Errorf("%s contains invalid control characters", field)
		}
	}
	return nil
}

func validateCronExpression(expr string) error {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil
	}
	if len(expr) > 100 {
		return fmt.Errorf("Cron 表达式过长")
	}
	parts := strings.Fields(expr)
	if len(parts) != 5 && len(parts) != 6 {
		return fmt.Errorf("Cron 表达式格式不正确，应为 5 位或 6 位格式")
	}
	if len(parts) == 5 {
		expr = "0 " + expr
	}
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, err := parser.Parse(expr); err != nil {
		return fmt.Errorf("Cron 表达式无效: %w", err)
	}
	return nil
}
