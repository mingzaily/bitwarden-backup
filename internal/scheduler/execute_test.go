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
