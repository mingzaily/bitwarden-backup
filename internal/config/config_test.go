package config

import (
	"strings"
	"testing"
)

func TestConfigValidateAdminPasswordMinimumLength(t *testing.T) {
	t.Setenv("AUTH_COOKIE_SECURE", "")

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "empty password", password: "", wantErr: true},
		{name: "seven characters", password: "1234567", wantErr: true},
		{name: "eight characters", password: "12345678", wantErr: false},
		{name: "eight unicode runes", password: "密码密码密码密码", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				ServerPort:    "8080",
				DBPath:        "./data/test.db",
				AdminPassword: tt.password,
			}

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "seven characters" && (err == nil || !strings.Contains(err.Error(), "8")) {
				t.Fatalf("Validate() error = %v, want message to mention minimum length 8", err)
			}
		})
	}
}
