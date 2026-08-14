package model

import "testing"

func TestIsOfficialServerURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{name: "global official endpoint", url: "https://vault.bitwarden.com", want: true},
		{name: "eu official endpoint with trailing slash", url: "HTTPS://VAULT.BITWARDEN.EU/", want: true},
		{name: "http is not trusted", url: "http://vault.bitwarden.com", want: false},
		{name: "lookalike hostname is not trusted", url: "https://vault.bitwarden.com.example.com", want: false},
		{name: "self hosted endpoint", url: "https://vault.example.com", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOfficialServerURL(tt.url); got != tt.want {
				t.Fatalf("IsOfficialServerURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}
