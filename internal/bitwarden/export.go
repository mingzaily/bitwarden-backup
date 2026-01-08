package bitwarden

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Unlock 解锁密码库
func (c *Client) Unlock(masterPassword string) error {
	if status, err := c.Status(); err == nil {
		if status == "unauthenticated" {
			return fmt.Errorf("bw status is unauthenticated; login required before unlock")
		}
	}

	stdin := masterPassword
	if !strings.HasSuffix(stdin, "\n") {
		stdin += "\n"
	}

	res, err := c.runBW([]string{"unlock", "--raw"}, stdin, nil)
	stderr := strings.TrimSpace(res.Stderr)
	if stderr != "" {
		log.Printf("[bitwarden] bw unlock stderr: %s", stderr)
	}
	if err != nil {
		stdout := strings.TrimSpace(res.Stdout)
		if stdout != "" {
			log.Printf("[bitwarden] bw unlock stdout: %s", stdout)
		}
		return fmt.Errorf("unlock failed (exit=%d): %w", res.ExitCode, err)
	}

	token := strings.TrimSpace(res.Stdout)
	if token == "" {
		stdout := strings.TrimSpace(res.Stdout)
		if stdout != "" {
			log.Printf("[bitwarden] bw unlock stdout: %s", stdout)
		}

		// 二次预检：若 vault 已是 unlocked，则允许继续（无需 session）
		if status, serr := c.Status(); serr == nil && status == "unlocked" {
			c.vaultUnlocked = true
			c.sessionToken = ""
			return nil
		}

		return fmt.Errorf("unlock returned empty session token (exit=%d); stdout=%s stderr=%s", res.ExitCode, strings.TrimSpace(res.Stdout), strings.TrimSpace(res.Stderr))
	}

	c.sessionToken = token
	c.vaultUnlocked = true
	log.Printf("[bitwarden] bw unlock ok (session token length=%d)", len(token))
	return nil
}

// Export 导出密码库数据
// password 参数为可选，仅在 format 为 "encrypted_json" 时需要提供
func (c *Client) Export(outputPath, format string, password ...string) error {
	if c.sessionToken == "" && !c.vaultUnlocked {
		return fmt.Errorf("vault is not unlocked, please unlock first")
	}

	// 确保输出目录存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 构建命令参数
	args := []string{"export", "--output", outputPath, "--format", format}
	if c.sessionToken != "" {
		args = append(args, "--session", c.sessionToken)
	}

	// 如果提供了密码参数，添加 --password 选项
	if len(password) > 0 && password[0] != "" {
		args = append(args, "--password", password[0])
	}

	res, err := c.runBW(args, "", nil)
	if err != nil {
		if strings.TrimSpace(res.Stdout) != "" {
			log.Printf("[bitwarden] bw export stdout: %s", strings.TrimSpace(res.Stdout))
		}
		if strings.TrimSpace(res.Stderr) != "" {
			log.Printf("[bitwarden] bw export stderr: %s", strings.TrimSpace(res.Stderr))
		}
		return fmt.Errorf("export failed (exit=%d): %w", res.ExitCode, err)
	}

	return nil
}
