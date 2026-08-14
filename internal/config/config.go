package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort       string
	DBPath           string
	LogLevel         string
	AppEnv           string
	AdminPassword    string
	AuthCookieSecure bool
}

func Load() *Config {
	return &Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		DBPath:           getEnv("DB_PATH", "./data/bitwarden-backup.db"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		AppEnv:           getEnv("APP_ENV", "production"),
		AdminPassword:    getEnv("BITWARDEN_BACKUP_ADMIN_PASSWORD", ""),
		AuthCookieSecure: getEnvAsBool("AUTH_COOKIE_SECURE", false),
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.AdminPassword) == "" {
		return fmt.Errorf("BITWARDEN_BACKUP_ADMIN_PASSWORD must be set")
	}
	if len([]rune(c.AdminPassword)) < 8 {
		return fmt.Errorf("BITWARDEN_BACKUP_ADMIN_PASSWORD must contain at least 8 characters")
	}
	if strings.TrimSpace(c.ServerPort) == "" {
		return fmt.Errorf("SERVER_PORT must not be empty")
	}
	if strings.TrimSpace(c.DBPath) == "" {
		return fmt.Errorf("DB_PATH must not be empty")
	}
	if raw := os.Getenv("AUTH_COOKIE_SECURE"); raw != "" {
		if _, err := strconv.ParseBool(raw); err != nil {
			return fmt.Errorf("AUTH_COOKIE_SECURE must be true or false")
		}
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
