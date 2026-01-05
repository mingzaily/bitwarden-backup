package bitwarden

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Unlock 解锁密码库
func (c *Client) Unlock(masterPassword string) error {
	cmd := exec.Command("bw", "unlock", masterPassword, "--raw")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unlock failed: %s, %w", string(output), err)
	}

	c.sessionToken = strings.TrimSpace(string(output))

	// 检查 session token 是否为空
	if c.sessionToken == "" {
		return fmt.Errorf("unlock returned empty session token, output: %s", string(output))
	}

	return nil
}

// Export 导出密码库数据
func (c *Client) Export(outputPath, format string) error {
	if c.sessionToken == "" {
		return fmt.Errorf("session token is empty, please unlock first")
	}

	// 确保输出目录存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	cmd := exec.Command("bw", "export", "--output", outputPath, "--format", format, "--session", c.sessionToken)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("export failed: %s, %w", string(output), err)
	}

	return nil
}
