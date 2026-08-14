package bitwarden

import "testing"

func TestAddLogDoesNotMarkExpectedLogoutAsError(t *testing.T) {
	client := NewClient()
	client.AddLog("bw logout (exit=1, 651ms)")
	client.AddLog("bw logout: already logged out")

	logs := client.GetLogs()
	if len(logs) != 2 {
		t.Fatalf("log count = %d, want 2", len(logs))
	}
	for _, entry := range logs {
		if entry.Level != "info" {
			t.Errorf("log %q level = %q, want info", entry.Message, entry.Level)
		}
	}
}

func TestNestedClientKeepsExpectedLogoutInformational(t *testing.T) {
	parent := NewClient()
	child := NewClientWithLogSink("server", parent.AddLogWithSource)
	child.AddLog("bw logout (exit=1, 651ms)")

	logs := parent.GetLogs()
	if len(logs) != 1 {
		t.Fatalf("parent log count = %d, want 1", len(logs))
	}
	if logs[0].Level != "info" {
		t.Fatalf("parent log level = %q, want info", logs[0].Level)
	}
}
