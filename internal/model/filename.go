package model

import (
	"fmt"
	"regexp"
	"strings"
)

// DefaultFilenameTemplate is used when a task does not specify a custom
// filename. The rendered timestamp is always YYYYMMDDHHmmss.
const DefaultFilenameTemplate = "bitwarden_encrypted_export_{time}.json"

var (
	filenamePlaceholderPattern = regexp.MustCompile(`\{([a-z_]+)\}`)
	backupTimestampPattern     = regexp.MustCompile(`^\d{14}$`)
)

// NormalizeFilenameTemplate keeps old tasks, which have no value for the
// newly added field, on the same safe default as newly created tasks.
func NormalizeFilenameTemplate(template string) string {
	if strings.TrimSpace(template) == "" {
		return DefaultFilenameTemplate
	}
	return strings.TrimSpace(template)
}

// ValidateFilenameTemplate validates the small, intentionally limited
// template language exposed by the task form. Empty values are accepted and
// normalized to DefaultFilenameTemplate by the request handler/scheduler.
func ValidateFilenameTemplate(template string) error {
	template = strings.TrimSpace(template)
	if template == "" {
		return nil
	}
	if len([]rune(template)) > 200 {
		return fmt.Errorf("filename_template is too long")
	}
	if !strings.Contains(template, "{time}") {
		return fmt.Errorf("filename_template must contain {time}")
	}
	if !strings.HasSuffix(strings.ToLower(template), ".json") {
		return fmt.Errorf("filename_template must end with .json")
	}
	if strings.ContainsAny(template, `/\\`) || strings.Contains(template, "..") {
		return fmt.Errorf("filename_template must be a single safe filename")
	}

	allowed := map[string]struct{}{
		"time":      {},
		"task_name": {},
		"medium":    {},
	}
	for _, match := range filenamePlaceholderPattern.FindAllStringSubmatch(template, -1) {
		if _, ok := allowed[match[1]]; !ok {
			return fmt.Errorf("filename_template contains unsupported placeholder {%s}", match[1])
		}
	}
	literal := filenamePlaceholderPattern.ReplaceAllString(template, "")
	if strings.ContainsAny(literal, "{}") {
		return fmt.Errorf("filename_template contains an invalid placeholder")
	}
	return nil
}

// ValidateBackupTimestamp keeps provider filenames stable and prevents a
// caller from accidentally rendering a timestamp with a different layout.
func ValidateBackupTimestamp(timestamp string) error {
	if !backupTimestampPattern.MatchString(timestamp) {
		return fmt.Errorf("backup timestamp must use YYYYMMDDHHmmss")
	}
	return nil
}

// RenderFilenameTemplate replaces the supported placeholders. The provider
// still runs the result through its filename safety helper before using it as
// a local or remote path component.
func RenderFilenameTemplate(template, taskName, timestamp, medium string) string {
	template = NormalizeFilenameTemplate(template)
	template = strings.ReplaceAll(template, "{time}", timestamp)
	template = strings.ReplaceAll(template, "{task_name}", taskName)
	template = strings.ReplaceAll(template, "{medium}", medium)
	return template
}
