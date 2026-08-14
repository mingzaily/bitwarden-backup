package model

import "testing"

func TestDestinationRequestClearsCredentialsWhenTypeChanges(t *testing.T) {
	destination := &BackupDestination{
		Type:               "s3",
		S3AccessKey:        "old-access-key",
		S3SecretKey:        "old-secret-key",
		EncryptionPassword: "old-encryption-password",
	}

	DestinationRequest{Type: "local", LocalPath: "/backups"}.ApplyTo(destination)

	if destination.S3AccessKey != "" || destination.S3SecretKey != "" || destination.EncryptionPassword != "" {
		t.Fatalf("old provider credentials were retained: %#v", destination)
	}
}

func TestDestinationRequestKeepsOmittedCredentialsOnSameType(t *testing.T) {
	destination := &BackupDestination{
		Type:               "webdav",
		WebDAVPassword:     "existing-password",
		EncryptionPassword: "existing-encryption-password",
		Encrypted:          true,
	}

	DestinationRequest{Type: "webdav", WebDAVURL: "https://dav.example.com", Encrypted: true}.ApplyTo(destination)

	if destination.WebDAVPassword != "existing-password" || destination.EncryptionPassword != "existing-encryption-password" {
		t.Fatalf("omitted credentials were not preserved: %#v", destination)
	}
}

func TestDestinationRequestClearsDisabledEncryptionPassword(t *testing.T) {
	destination := &BackupDestination{Type: "local", Encrypted: true, EncryptionPassword: "existing"}

	DestinationRequest{Type: "local", LocalPath: "/backups", Encrypted: false}.ApplyTo(destination)

	if destination.EncryptionPassword != "" {
		t.Fatal("encryption password should be cleared when encryption is disabled")
	}
}
