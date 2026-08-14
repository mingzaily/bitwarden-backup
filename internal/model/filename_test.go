package model

import "testing"

func TestFilenameTemplateDefaultsAndRenders(t *testing.T) {
	template := NormalizeFilenameTemplate("")
	if template != DefaultFilenameTemplate {
		t.Fatalf("default template = %q, want %q", template, DefaultFilenameTemplate)
	}
	if err := ValidateFilenameTemplate(template); err != nil {
		t.Fatalf("ValidateFilenameTemplate(default) returned error: %v", err)
	}
	got := RenderFilenameTemplate(template, "每日备份", "20251204092928", "")
	if got != "bitwarden_encrypted_export_20251204092928.json" {
		t.Fatalf("RenderFilenameTemplate() = %q", got)
	}
}

func TestFilenameTemplateSupportsTaskAndMedium(t *testing.T) {
	template := "{task_name}_{medium}_{time}.json"
	if err := ValidateFilenameTemplate(template); err != nil {
		t.Fatalf("ValidateFilenameTemplate() returned error: %v", err)
	}
	got := RenderFilenameTemplate(template, "nightly", "20251204092928", "webdav")
	if got != "nightly_webdav_20251204092928.json" {
		t.Fatalf("RenderFilenameTemplate() = %q", got)
	}
}

func TestFilenameTemplateRejectsUnsafeOrIncompleteValues(t *testing.T) {
	for _, template := range []string{
		"backup.json",
		"backup_{time}.zip",
		"../backup_{time}.json",
		"backup_{unknown}_{time}.json",
		"backup_{time}.json}",
	} {
		if err := ValidateFilenameTemplate(template); err == nil {
			t.Errorf("ValidateFilenameTemplate(%q) returned nil", template)
		}
	}
}

func TestValidateBackupTimestamp(t *testing.T) {
	if err := ValidateBackupTimestamp("20251204092928"); err != nil {
		t.Fatalf("valid timestamp rejected: %v", err)
	}
	for _, timestamp := range []string{"20251204_092928", "2025120409292", "202512040929280"} {
		if err := ValidateBackupTimestamp(timestamp); err == nil {
			t.Errorf("timestamp %q was accepted", timestamp)
		}
	}
}
