package bitwarden

import (
	"reflect"
	"testing"
)

func TestBuildExportArgsUsesExportPasswordFlag(t *testing.T) {
	got := buildExportArgs("/tmp/backup.json", "encrypted_json", "session-token", "export-password")
	want := []string{
		"export", "--output", "/tmp/backup.json", "--format", "encrypted_json",
		"--session", "session-token", "--password", "export-password",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildExportArgs() = %#v, want %#v", got, want)
	}
}

func TestBuildExportArgsOmitsEmptyPassword(t *testing.T) {
	got := buildExportArgs("/tmp/backup.json", "json", "session-token", "")
	want := []string{"export", "--output", "/tmp/backup.json", "--format", "json", "--session", "session-token"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildExportArgs() = %#v, want %#v", got, want)
	}
}

func TestRedactBWArgsHidesPasswordAndSession(t *testing.T) {
	got := redactBWArgs([]string{"export", "--session", "session-token", "--password", "export-password"})
	want := []string{"export", "--session", "***", "--password", "***"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("redactBWArgs() = %#v, want %#v", got, want)
	}
}
