package bitwarden

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// passwordPromptRegex 匹配密码提示信息
var passwordPromptRegex = regexp.MustCompile(`\?\s*Master password:\s*(\[input is hidden\]|\[hidden\])\s*`)

// cleanPasswordPrompts 清理 stderr 中的密码提示信息
func cleanPasswordPrompts(s string) string {
	// 移除所有密码提示
	cleaned := passwordPromptRegex.ReplaceAllString(s, "")
	// 清理多余的空白
	cleaned = strings.TrimSpace(cleaned)
	// 如果清理后只剩空白，返回空字符串
	if cleaned == "" {
		return ""
	}
	return cleaned
}

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
	// 清理 stderr 中的密码提示信息，只保留有意义的错误信息
	if stderr != "" {
		// 移除重复的密码提示 "? Master password: [input is hidden]" 或 "[hidden]"
		cleanStderr := sanitizeBWOutput(stderr)
		if cleanStderr != "" {
			c.AddLog(fmt.Sprintf("bw unlock stderr: %s", cleanStderr))
		}
	}
	if err != nil {
		stdout := strings.TrimSpace(res.Stdout)
		if stdout != "" {
			c.AddLog(fmt.Sprintf("bw unlock stdout: %s", stdout))
		}
		// 检测登录状态损坏的情况
		if strings.Contains(stderr, "not logged in") || strings.Contains(stderr, "You are not logged in") {
			return &ErrNotLoggedIn{Msg: fmt.Sprintf("unlock failed: %s", stderr)}
		}
		return fmt.Errorf("unlock failed (exit=%d): %w", res.ExitCode, err)
	}

	token := strings.TrimSpace(res.Stdout)
	if token == "" {
		stdout := strings.TrimSpace(res.Stdout)
		if stdout != "" {
			c.AddLog(fmt.Sprintf("bw unlock stdout: %s", stdout))
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
	c.AddLog(fmt.Sprintf("bw unlock ok (session token length=%d)", len(token)))
	return nil
}

// ErrNotLoggedIn 表示登录状态已失效，需要重新登录
type ErrNotLoggedIn struct {
	Msg string
}

func (e *ErrNotLoggedIn) Error() string {
	return e.Msg
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
			c.AddLog(fmt.Sprintf("bw export stdout: %s", strings.TrimSpace(res.Stdout)))
		}
		if strings.TrimSpace(res.Stderr) != "" {
			c.AddLog(fmt.Sprintf("bw export stderr: %s", strings.TrimSpace(res.Stderr)))
		}
		return fmt.Errorf("export failed (exit=%d): %w", res.ExitCode, err)
	}

	return nil
}
