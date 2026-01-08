package scheduler

import (
	"fmt"

	"github.com/mingzaily/bitwarden-backup/internal/model"
)

func (s *Scheduler) backupToDestination(dest model.BackupDestination, sourceFile, taskName, timestamp string) error {
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
