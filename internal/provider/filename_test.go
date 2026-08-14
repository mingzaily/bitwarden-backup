package provider

import (
	"testing"

	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func TestRenderBackupFilenameUsesDestinationMedium(t *testing.T) {
	for _, test := range []struct {
		providerType string
		want         string
	}{
		{providerType: "local", want: "nightly_local_20251204092928.json"},
		{providerType: "webdav", want: "nightly_webdav_20251204092928.json"},
		{providerType: "s3", want: "nightly_oss_20251204092928.json"},
	} {
		got := renderBackupFilename(BackupContext{
			TaskName:         "nightly",
			Timestamp:        "20251204092928",
			FilenameTemplate: "{task_name}_{medium}_{time}.json",
			Destination:      model.BackupDestination{Type: test.providerType},
		})
		if got != test.want {
			t.Errorf("renderBackupFilename(%q) = %q, want %q", test.providerType, got, test.want)
		}
	}
	serverDefault := renderBackupFilename(BackupContext{
		TaskName:         "nightly",
		Timestamp:        "20251204092928",
		FilenameTemplate: model.DefaultFilenameTemplate,
		Destination:      model.BackupDestination{Type: "server"},
	})
	if serverDefault != "bitwarden_encrypted_export_20251204092928.json" {
		t.Errorf("renderBackupFilename(server) = %q", serverDefault)
	}
}

func TestMatchesBackupFilenameScopesRetentionToTaskTemplate(t *testing.T) {
	ctx := BackupContext{
		TaskName:         "nightly",
		Timestamp:        "20251204092928",
		FilenameTemplate: "{task_name}_{medium}_{time}.json",
		Destination:      model.BackupDestination{Type: "webdav"},
	}

	for _, name := range []string{
		"nightly_webdav_20251204092928.json",
		"nightly_webdav_20991231235959.json",
	} {
		if !matchesBackupFilename(name, ctx) {
			t.Errorf("matchesBackupFilename(%q) = false", name)
		}
	}
	for _, name := range []string{
		"bitwarden_encrypted_export_20251204092928.json",
		"other_webdav_20251204092928.json",
		"notes_20251204092928.json",
		"backup.txt",
	} {
		if matchesBackupFilename(name, ctx) {
			t.Errorf("matchesBackupFilename(%q) = true", name)
		}
	}
}

func TestMatchesBackupFilenameKeepsLegacyTaskFiles(t *testing.T) {
	ctx := BackupContext{
		TaskName:         "nightly",
		FilenameTemplate: model.DefaultFilenameTemplate,
		Destination:      model.BackupDestination{Type: "local"},
	}
	if !matchesBackupFilename("backup_nightly_20251204_092928.000000000.json", ctx) {
		t.Fatal("legacy task filename was not recognized")
	}
}
