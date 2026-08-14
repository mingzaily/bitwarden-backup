package scheduler

import (
	"context"
	"errors"
	"testing"

	"github.com/mingzaily/bitwarden-backup/internal/model"
	"github.com/mingzaily/bitwarden-backup/internal/provider"
)

type cleanupFailureProvider struct{}

func (p *cleanupFailureProvider) Type() string { return "test-cleanup-failure" }

func (p *cleanupFailureProvider) Backup(provider.BackupContext) (string, error) {
	return "/backups/backup.json", nil
}

func (p *cleanupFailureProvider) Cleanup(provider.BackupContext, int) (int, error) {
	return 0, errors.New("retention backend unavailable")
}

func TestBackupToDestinationReturnsCleanupErrorWithArtifact(t *testing.T) {
	provider.GetRegistry().Register(&cleanupFailureProvider{})

	path, err := (&Scheduler{}).backupToDestination(
		context.Background(),
		model.BackupDestination{
			Name:           "test destination",
			Type:           "test-cleanup-failure",
			MaxBackupCount: 1,
		},
		"/tmp/source.json",
		"test task",
		"20251204092928",
		model.DefaultFilenameTemplate,
		nil,
	)
	if err == nil {
		t.Fatal("cleanup failure should fail the destination")
	}
	if path != "/backups/backup.json" {
		t.Fatalf("backup path = %q, want uploaded artifact path", path)
	}
}
