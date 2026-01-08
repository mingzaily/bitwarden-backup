package bitwarden

import (
	"fmt"
	"log"
	"strings"
)

// Import 导入数据到密码库
func (c *Client) Import(inputPath, format string) error {
	if c.sessionToken == "" && !c.vaultUnlocked {
		return fmt.Errorf("vault is not unlocked, please unlock first")
	}

	args := []string{"import", format, inputPath}
	if c.sessionToken != "" {
		args = append(args, "--session", c.sessionToken)
	}

	res, err := c.runBW(args, "", nil)
	if err != nil {
		if strings.TrimSpace(res.Stdout) != "" {
			log.Printf("[bitwarden] bw import stdout: %s", strings.TrimSpace(res.Stdout))
		}
		if strings.TrimSpace(res.Stderr) != "" {
			log.Printf("[bitwarden] bw import stderr: %s", strings.TrimSpace(res.Stderr))
		}
		return fmt.Errorf("import failed (exit=%d): %w", res.ExitCode, err)
	}

	return nil
}

// Logout 登出
func (c *Client) Logout() error {
	res, err := c.runBW([]string{"logout"}, "", nil)
	if err != nil {
		if strings.TrimSpace(res.Stdout) != "" {
			log.Printf("[bitwarden] bw logout stdout: %s", strings.TrimSpace(res.Stdout))
		}
		if strings.TrimSpace(res.Stderr) != "" {
			log.Printf("[bitwarden] bw logout stderr: %s", strings.TrimSpace(res.Stderr))
		}
		return fmt.Errorf("logout failed (exit=%d): %w", res.ExitCode, err)
	}

	c.sessionToken = ""
	c.vaultUnlocked = false
	return nil
}
