package scheduler

import (
	"fmt"

	"github.com/mingzaily/bitwarden-backup/internal/database"
)

// backupToDestination 备份到单个目标
func (s *Scheduler) backupToDestination(dest database.BackupDestination, sourceFile, taskName, timestamp string) error {
	switch dest.Type {
	case "local":
		return s.backupToLocal(dest, sourceFile, taskName, timestamp)
	case "webdav":
		return s.backupToWebDAV(dest, sourceFile, taskName, timestamp)
	case "server":
		return s.backupToServer(dest, sourceFile)
	default:
		return fmt.Errorf("unknown destination type: %s", dest.Type)
	}
}
