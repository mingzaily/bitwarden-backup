package bitwarden

import (
	"fmt"
	"os/exec"
)

// Import 导入数据到密码库
func (c *Client) Import(inputPath, format string) error {
	if c.sessionToken == "" {
		return fmt.Errorf("session token is empty, please unlock first")
	}

	cmd := exec.Command("bw", "import", format, inputPath, "--session", c.sessionToken)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("import failed: %s, %w", string(output), err)
	}

	return nil
}

// Logout 登出
func (c *Client) Logout() error {
	cmd := exec.Command("bw", "logout")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("logout failed: %s, %w", string(output), err)
	}

	c.sessionToken = ""
	return nil
}
