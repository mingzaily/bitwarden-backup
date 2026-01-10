package bitwarden

import (
	"context"
	"fmt"
	"strings"
)

// Import 导入数据到密码库
func (c *Client) Import(ctx context.Context, inputPath, format string) error {
	if c.sessionToken == "" && !c.vaultUnlocked {
		return fmt.Errorf("vault is not unlocked, please unlock first")
	}

	args := []string{"import", format, inputPath}
	if c.sessionToken != "" {
		args = append(args, "--session", c.sessionToken)
	}

	res, err := c.runBW(ctx, args, "", nil)
	if err != nil {
		if strings.TrimSpace(res.Stdout) != "" {
			c.AddLog(fmt.Sprintf("bw import stdout: %s", strings.TrimSpace(res.Stdout)))
		}
		if strings.TrimSpace(res.Stderr) != "" {
			c.AddLog(fmt.Sprintf("bw import stderr: %s", strings.TrimSpace(res.Stderr)))
		}
		return fmt.Errorf("import failed (exit=%d): %w", res.ExitCode, err)
	}

	return nil
}

// Logout 登出
func (c *Client) Logout(ctx context.Context) error {
	res, err := c.runBW(ctx, []string{"logout"}, "", nil)
	if err != nil {
		stderr := strings.TrimSpace(res.Stderr)
		// "You are not logged in" 意味着已经是登出状态，不算错误
		if strings.Contains(stderr, "You are not logged in") {
			c.AddLog("bw logout: already logged out")
			c.sessionToken = ""
			c.vaultUnlocked = false
			return nil
		}
		if strings.TrimSpace(res.Stdout) != "" {
			c.AddLog(fmt.Sprintf("bw logout stdout: %s", strings.TrimSpace(res.Stdout)))
		}
		if stderr != "" {
			c.AddLog(fmt.Sprintf("bw logout stderr: %s", stderr))
		}
		return fmt.Errorf("logout failed (exit=%d): %w", res.ExitCode, err)
	}

	c.sessionToken = ""
	c.vaultUnlocked = false
	return nil
}
