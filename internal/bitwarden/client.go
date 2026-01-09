package bitwarden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// LogEntry 单条执行日志
type LogEntry struct {
	Time    string `json:"time"`
	Message string `json:"message"`
}

// Client Bitwarden CLI 客户端
type Client struct {
	sessionToken  string
	serverURL     string
	vaultUnlocked bool
	logs          []LogEntry
}

// NewClient 创建新的 Bitwarden 客户端
func NewClient() *Client {
	return &Client{logs: make([]LogEntry, 0)}
}

// ansiRegex 匹配 ANSI 转义序列
// 格式: ESC[ + 参数 + 命令字母，其中 ESC 可以是 \x1b 或 \033
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// 敏感输出正则表达式（匹配密码提示、输入隐藏等）
var sensitiveOutputRegex = regexp.MustCompile(`(?im)(?:master\s*password|password:|input\s+is\s+hidden|\[hidden\]|\[input\s+is\s+hidden\]).*`)
// Session token 正则表达式（匹配长字符串）
var tokenRegex = regexp.MustCompile(`[a-zA-Z0-9+/]{64,}`)

// stripANSI 移除字符串中的 ANSI 转义序列
func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// sanitizeBWOutput 统一脱敏 Bitwarden CLI 输出
func sanitizeBWOutput(s string) string {
	s = stripANSI(s)
	// 移除包含敏感关键词的整行
	lines := strings.Split(s, "\n")
	var cleaned []string
	for _, line := range lines {
		if !sensitiveOutputRegex.MatchString(line) {
			cleaned = append(cleaned, line)
		}
	}
	s = strings.Join(cleaned, "\n")
	// 掩码可能的 token
	s = tokenRegex.ReplaceAllString(s, "***")
	return strings.TrimSpace(s)
}

// AddLog 添加一条日志
func (c *Client) AddLog(message string) {
	// 统一脱敏处理
	cleanMessage := sanitizeBWOutput(message)
	if cleanMessage == "" {
		return // 如果脱敏后为空，不记录日志
	}
	c.logs = append(c.logs, LogEntry{
		Time:    time.Now().Format("2006/01/02 15:04:05"),
		Message: cleanMessage,
	})
	log.Printf("[bitwarden] %s", cleanMessage)
}

// GetLogs 获取所有日志
func (c *Client) GetLogs() []LogEntry {
	return c.logs
}

// ClearLogs 清空日志
func (c *Client) ClearLogs() {
	c.logs = make([]LogEntry, 0)
}

type bwExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func redactBWArgs(args []string) []string {
	redacted := make([]string, len(args))
	copy(redacted, args)

	for i := 0; i < len(redacted); i++ {
		if (redacted[i] == "--session" || redacted[i] == "--password") && i+1 < len(redacted) {
			redacted[i+1] = "***"
			i++
			continue
		}
		if strings.HasPrefix(redacted[i], "--session=") {
			redacted[i] = "--session=***"
			continue
		}
		if strings.HasPrefix(redacted[i], "--password=") {
			redacted[i] = "--password=***"
			continue
		}
	}

	return redacted
}

func (c *Client) runBW(args []string, stdin string, extraEnv map[string]string) (bwExecResult, error) {
	cmd := exec.Command("bw", args...)
	cmd.Env = os.Environ()
	for k, v := range extraEnv {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	} else if err != nil {
		exitCode = -1
	}

	res := bwExecResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: exitCode,
	}

	c.AddLog(fmt.Sprintf("bw %s (exit=%d)", strings.Join(redactBWArgs(args), " "), exitCode))
	return res, err
}

type bwStatusResponse struct {
	Status string `json:"status"`
}

// Status 获取当前 vault 状态（unauthenticated / locked / unlocked）
func (c *Client) Status() (string, error) {
	res, err := c.runBW([]string{"status"}, "", nil)
	stdout := strings.TrimSpace(res.Stdout)
	stderr := strings.TrimSpace(res.Stderr)

	if stdout != "" {
		c.AddLog(fmt.Sprintf("bw status stdout: %s", stdout))
	}
	if stderr != "" {
		c.AddLog(fmt.Sprintf("bw status stderr: %s", stderr))
	}
	if err != nil {
		return "", fmt.Errorf("bw status failed (exit=%d)", res.ExitCode)
	}

	var parsed bwStatusResponse
	if uerr := json.Unmarshal([]byte(stdout), &parsed); uerr != nil {
		return "", fmt.Errorf("failed to parse bw status output: %w", uerr)
	}
	if parsed.Status == "" {
		return "", fmt.Errorf("bw status returned empty status")
	}
	return parsed.Status, nil
}

// ConfigServer 配置服务器地址
func (c *Client) ConfigServer(serverURL string) error {
	res, err := c.runBW([]string{"config", "server", serverURL}, "", nil)
	if err != nil {
		if strings.TrimSpace(res.Stdout) != "" {
			c.AddLog(fmt.Sprintf("bw config server stdout: %s", strings.TrimSpace(res.Stdout)))
		}
		if strings.TrimSpace(res.Stderr) != "" {
			c.AddLog(fmt.Sprintf("bw config server stderr: %s", strings.TrimSpace(res.Stderr)))
		}
		return fmt.Errorf("config server failed (exit=%d): %w", res.ExitCode, err)
	}
	c.serverURL = serverURL
	return nil
}

// Login 登录到 Bitwarden
func (c *Client) Login(clientID, clientSecret string) error {
	res, err := c.runBW(
		[]string{"login", "--apikey"},
		"",
		map[string]string{
			"BW_CLIENTID":     clientID,
			"BW_CLIENTSECRET": clientSecret,
		},
	)
	if err != nil {
		if strings.TrimSpace(res.Stdout) != "" {
			c.AddLog(fmt.Sprintf("bw login stdout: %s", strings.TrimSpace(res.Stdout)))
		}
		if strings.TrimSpace(res.Stderr) != "" {
			c.AddLog(fmt.Sprintf("bw login stderr: %s", strings.TrimSpace(res.Stderr)))
		}
		return fmt.Errorf("login failed (exit=%d): %w", res.ExitCode, err)
	}
	return nil
}
