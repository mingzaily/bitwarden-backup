package webdav

import (
	"strings"
)

// Client WebDAV 客户端
type Client struct {
	baseURL  string
	username string
	password string
}

// NewClient 创建 WebDAV 客户端
func NewClient(baseURL, username, password string) *Client {
	return &Client{
		baseURL:  strings.TrimSuffix(baseURL, "/"),
		username: username,
		password: password,
	}
}
