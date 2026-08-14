package webdav

import "testing"

func TestRequestURLJoinsBaseAndRemotePath(t *testing.T) {
	client := NewClient("https://dav.example.test/root/", "user", "password")

	got, err := client.requestURL("/vault/backup.zip")
	if err != nil {
		t.Fatalf("requestURL returned error: %v", err)
	}
	if want := "https://dav.example.test/root/vault/backup.zip"; got != want {
		t.Fatalf("requestURL() = %q, want %q", got, want)
	}
}

func TestRequestURLRejectsParentTraversal(t *testing.T) {
	client := NewClient("https://dav.example.test/root", "user", "password")

	for _, remotePath := range []string{"../backup.zip", "/vault/../backup.zip", "vault\\backup.zip"} {
		if _, err := client.requestURL(remotePath); err == nil {
			t.Errorf("requestURL(%q) accepted an unsafe path", remotePath)
		}
	}
}

func TestRequestURLRejectsInsecureRemoteBase(t *testing.T) {
	client := NewClient("http://dav.example.test/root", "user", "password")
	if _, err := client.requestURL("backup.zip"); err == nil {
		t.Fatal("requestURL accepted a non-loopback HTTP base URL")
	}
}
