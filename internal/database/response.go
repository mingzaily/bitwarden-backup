package database

// ServerConfigResponse 服务器配置响应（隐藏敏感信息）
type ServerConfigResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	ServerURL      string `json:"server_url"`
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	MasterPassword string `json:"master_password"`
	IsOfficial     bool   `json:"is_official"`
	Enabled        bool   `json:"enabled"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// ToResponse 转换为响应格式（隐藏敏感信息）
func (s *ServerConfig) ToResponse() ServerConfigResponse {
	return ServerConfigResponse{
		ID:             s.ID,
		Name:           s.Name,
		ServerURL:      s.ServerURL,
		ClientID:       maskSensitive(s.ClientID),
		ClientSecret:   maskSensitive(s.ClientSecret),
		MasterPassword: maskSensitive(s.MasterPassword),
		IsOfficial:     s.IsOfficial,
		Enabled:        s.Enabled,
		CreatedAt:      s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      s.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// BackupDestinationResponse 备份目标响应（隐藏敏感信息）
type BackupDestinationResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	TypeLabel      string `json:"type_label"`
	DisplayPath    string `json:"display_path"`
	LocalPath      string `json:"local_path"`
	WebDAVURL      string `json:"webdav_url"`
	WebDAVUsername string `json:"webdav_username"`
	WebDAVPassword string `json:"webdav_password"`
	WebDAVPath         string `json:"webdav_path"`
	TargetServerID     *uint  `json:"target_server_id"`
	Encrypted          bool   `json:"encrypted"`
	EncryptionPassword string `json:"encryption_password"`
	Enabled            bool   `json:"enabled"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// ToResponse 转换为响应格式（隐藏敏感信息）
func (d *BackupDestination) ToResponse() BackupDestinationResponse {
	// 生成类型标签和展示路径
	typeLabel, displayPath := d.getTypeInfo()

	return BackupDestinationResponse{
		ID:                 d.ID,
		Name:               d.Name,
		Type:               d.Type,
		TypeLabel:          typeLabel,
		DisplayPath:        displayPath,
		LocalPath:          d.LocalPath,
		WebDAVURL:          d.WebDAVURL,
		WebDAVUsername:     d.WebDAVUsername,
		WebDAVPassword:     maskSensitive(d.WebDAVPassword),
		WebDAVPath:         d.WebDAVPath,
		TargetServerID:     d.TargetServerID,
		Encrypted:          d.Encrypted,
		EncryptionPassword: maskSensitive(d.EncryptionPassword),
		Enabled:            d.Enabled,
		CreatedAt:          d.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:          d.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// maskSensitive 隐藏敏感信息，只显示前后几个字符
func maskSensitive(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

// getTypeInfo 获取类型标签和展示路径
func (d *BackupDestination) getTypeInfo() (string, string) {
	switch d.Type {
	case "local":
		return "本地存储", d.LocalPath
	case "webdav":
		path := d.WebDAVURL
		if d.WebDAVPath != "" {
			path += d.WebDAVPath
		}
		return "WebDAV", path
	case "s3":
		path := "s3://" + d.S3Bucket
		if d.S3Path != "" {
			path += d.S3Path
		}
		return "S3 对象存储", path
	case "server":
		return "目标服务器", "N/A"
	default:
		return d.Type, "N/A"
	}
}
