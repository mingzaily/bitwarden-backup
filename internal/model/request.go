package model

// EnabledRequest is used by the explicit status endpoints. A pointer makes
// an omitted field distinguishable from an intentional false value.
type EnabledRequest struct {
	Enabled *bool `json:"enabled"`
}

// DestinationRequest is the API representation used to create or update a
// backup destination. It deliberately excludes database IDs, timestamps and
// preloaded associations from the request surface.
type DestinationRequest struct {
	Name               string `json:"name"`
	Type               string `json:"type"`
	LocalPath          string `json:"local_path"`
	WebDAVURL          string `json:"webdav_url"`
	WebDAVUsername     string `json:"webdav_username"`
	WebDAVPassword     string `json:"webdav_password"`
	WebDAVPath         string `json:"webdav_path"`
	S3Endpoint         string `json:"s3_endpoint"`
	S3Region           string `json:"s3_region"`
	S3Bucket           string `json:"s3_bucket"`
	S3AccessKey        string `json:"s3_access_key"`
	S3SecretKey        string `json:"s3_secret_key"`
	S3Path             string `json:"s3_path"`
	TargetServerID     *uint  `json:"target_server_id"`
	Encrypted          bool   `json:"encrypted"`
	EncryptionPassword string `json:"encryption_password"`
	MaxBackupCount     int    `json:"max_backup_count"`
	Enabled            *bool  `json:"enabled"`
}

// ToDestination converts a create request into a persistence model.
func (r DestinationRequest) ToDestination() *BackupDestination {
	enabled := true
	if r.Enabled != nil {
		enabled = *r.Enabled
	}

	destination := &BackupDestination{Enabled: enabled}
	r.applyTo(destination)
	return destination
}

// ApplyTo copies updateable fields onto an existing destination. Empty secret
// values intentionally keep the stored secret unchanged during an update.
func (r DestinationRequest) ApplyTo(destination *BackupDestination) {
	r.applyTo(destination)
	if r.Enabled != nil {
		destination.Enabled = *r.Enabled
	}
}

func (r DestinationRequest) applyTo(destination *BackupDestination) {
	previousType := destination.Type
	if previousType != "" && previousType != r.Type {
		// A type change must not retain credentials for the old provider.
		destination.WebDAVPassword = ""
		destination.S3AccessKey = ""
		destination.S3SecretKey = ""
		destination.EncryptionPassword = ""
	}

	destination.Name = r.Name
	destination.Type = r.Type
	destination.LocalPath = r.LocalPath
	destination.WebDAVURL = r.WebDAVURL
	destination.WebDAVUsername = r.WebDAVUsername
	destination.WebDAVPath = r.WebDAVPath
	destination.S3Endpoint = r.S3Endpoint
	destination.S3Region = r.S3Region
	destination.S3Bucket = r.S3Bucket
	destination.S3Path = r.S3Path
	destination.TargetServerID = r.TargetServerID
	destination.Encrypted = r.Encrypted
	destination.MaxBackupCount = r.MaxBackupCount
	if !r.Encrypted {
		destination.EncryptionPassword = ""
	}

	if r.WebDAVPassword != "" {
		destination.WebDAVPassword = r.WebDAVPassword
	}
	if r.S3AccessKey != "" {
		destination.S3AccessKey = r.S3AccessKey
	}
	if r.S3SecretKey != "" {
		destination.S3SecretKey = r.S3SecretKey
	}
	if r.EncryptionPassword != "" {
		destination.EncryptionPassword = r.EncryptionPassword
	}
}
