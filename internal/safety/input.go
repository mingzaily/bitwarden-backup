package safety

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"unicode"
)

// Filename converts user-controlled labels into a single path component.
// It deliberately keeps the set small so a task name can never introduce a
// separator, dot path, control character, or shell metacharacter.
func Filename(value string) string {
	value = strings.TrimSpace(value)
	var b strings.Builder
	lastUnderscore := false
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
			b.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore {
			b.WriteByte('_')
			lastUnderscore = true
		}
	}
	name := strings.Trim(b.String(), "._")
	if name == "" {
		return "task"
	}
	if len(name) > 80 {
		name = name[:80]
		name = strings.Trim(name, "._")
	}
	return name
}

func ValidateName(value, field string, max int) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	if len([]rune(value)) > max {
		return fmt.Errorf("%s is too long", field)
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return fmt.Errorf("%s contains invalid control characters", field)
		}
	}
	return nil
}

func ValidateURL(raw, field string, allowHTTP bool) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("%s is required", field)
	}
	if len([]rune(raw)) > 255 {
		return fmt.Errorf("%s is too long", field)
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("%s must be a valid URL", field)
	}
	if u.Scheme != "https" {
		if u.Scheme != "http" || !allowHTTP || !isLoopbackHost(u.Hostname()) {
			return fmt.Errorf("%s must use https (http is only allowed for localhost)", field)
		}
	}
	if strings.ContainsAny(raw, "\r\n") {
		return fmt.Errorf("%s contains invalid characters", field)
	}
	return nil
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func ValidateRemotePath(raw, field string) error {
	if len([]rune(raw)) > 255 {
		return fmt.Errorf("%s is too long", field)
	}
	if strings.ContainsAny(raw, "\r\n\\") {
		return fmt.Errorf("%s contains invalid path characters", field)
	}
	for _, segment := range strings.Split(strings.Trim(raw, "/"), "/") {
		if segment == ".." {
			return fmt.Errorf("%s must not contain parent path segments", field)
		}
	}
	return nil
}
