package bitwarden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Client Bitwarden CLI 客户端
type Client struct {
	sessionToken  string
	serverURL     string
	vaultUnlocked bool
}

// NewClient 创建新的 Bitwarden 客户端
func NewClient() *Client {
	return &Client{}
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

	log.Printf("[bitwarden] exec: bw %s (exit=%d)", strings.Join(redactBWArgs(args), " "), exitCode)
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
		log.Printf("[bitwarden] bw status stdout: %s", stdout)
	}
	if stderr != "" {
		log.Printf("[bitwarden] bw status stderr: %s", stderr)
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
			log.Printf("[bitwarden] bw config server stdout: %s", strings.TrimSpace(res.Stdout))
		}
		if strings.TrimSpace(res.Stderr) != "" {
			log.Printf("[bitwarden] bw config server stderr: %s", strings.TrimSpace(res.Stderr))
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
			log.Printf("[bitwarden] bw login stdout: %s", strings.TrimSpace(res.Stdout))
		}
		if strings.TrimSpace(res.Stderr) != "" {
			log.Printf("[bitwarden] bw login stderr: %s", strings.TrimSpace(res.Stderr))
		}
		return fmt.Errorf("login failed (exit=%d): %w", res.ExitCode, err)
	}
	return nil
}
