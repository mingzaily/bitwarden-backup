package scheduler

import "testing"

func TestTriggerTaskDeduplicatesQueuedTask(t *testing.T) {
	s := New()

	if !s.TriggerTask(42) {
		t.Fatal("first trigger should be accepted")
	}
	if s.TriggerTask(42) {
		t.Fatal("duplicate trigger should be rejected")
	}
	if got := <-s.taskQueue; got != 42 {
		t.Fatalf("queued task ID = %d, want 42", got)
	}
}

func TestNextBackupTimestampIsMonotonic(t *testing.T) {
	s := New()

	first := s.nextBackupTimestamp()
	second := s.nextBackupTimestamp()
	if second <= first {
		t.Fatalf("timestamps are not monotonic: first=%s second=%s", first, second)
	}
	if len(first) != 14 || len(second) != 14 {
		t.Fatalf("timestamps must use YYYYMMDDHHmmss format: first=%s second=%s", first, second)
	}
}
