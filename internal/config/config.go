package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	DBPath     string
	LogLevel   string
	AppEnv     string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "./data/bitwarden-backup.db"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		AppEnv:     getEnv("APP_ENV", "production"),
	}
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
