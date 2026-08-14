package safety

import "testing"

func TestFilenameRemovesPathSeparators(t *testing.T) {
	got := Filename("../../vault\nbackup:prod")
	if got != "vault_backup_prod" {
		t.Fatalf("Filename() = %q", got)
	}
}

func TestValidateRemotePathRejectsParentSegments(t *testing.T) {
	if err := ValidateRemotePath("/vault/../secrets", "path"); err == nil {
		t.Fatal("expected parent path segment to be rejected")
	}
}
