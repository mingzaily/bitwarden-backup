package bitwarden

import (
	"fmt"
	"os"
	"os/exec"
)

// Client Bitwarden CLI 客户端
type Client struct {
	sessionToken string
	serverURL    string
}

// NewClient 创建新的 Bitwarden 客户端
func NewClient() *Client {
	return &Client{}
}

// ConfigServer 配置服务器地址
func (c *Client) ConfigServer(serverURL string) error {
	cmd := exec.Command("bw", "config", "server", serverURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("config server failed: %s, %w", string(output), err)
	}
	c.serverURL = serverURL
	return nil
}

// Login 登录到 Bitwarden
func (c *Client) Login(clientID, clientSecret string) error {
	cmd := exec.Command("bw", "login", "--apikey")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("BW_CLIENTID=%s", clientID),
		fmt.Sprintf("BW_CLIENTSECRET=%s", clientSecret),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("login failed: %s, %w", string(output), err)
	}
	return nil
}
