package scheduler

import "time"

// nextBackupTimestamp returns a process-local monotonic timestamp. Backup
// filenames only carry second precision, so this prevents two destinations or
// manually triggered runs from producing the same filename within one second.
func (s *Scheduler) nextBackupTimestamp() string {
	now := time.Now().Truncate(time.Second)

	s.timestampMu.Lock()
	defer s.timestampMu.Unlock()

	if !s.lastTimestamp.IsZero() && !now.After(s.lastTimestamp) {
		now = s.lastTimestamp.Add(time.Second)
	}
	s.lastTimestamp = now

	return now.Format("20060102150405")
}
