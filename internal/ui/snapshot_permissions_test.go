package ui

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWorkerStallDumpDirAndFilesAreOwnerOnly(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	dir := writeWorkerStallDump(2*time.Minute, time.Now().Add(-time.Minute))
	if dir == "" {
		t.Fatal("no stall dump written")
	}

	for path, want := range map[string]os.FileMode{
		dir:                                  0o700,
		filepath.Dir(dir):                    0o700,
		filepath.Join(dir, "stall.json"):     0o600,
		filepath.Join(dir, "goroutines.txt"): 0o600,
	} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat %s: %v", path, err)
		}
		if got := info.Mode().Perm(); got != want {
			t.Errorf("%s mode = %o, want %o", path, got, want)
		}
	}
}

func TestTranscriptCopiedIntoASnapshotIsOwnerOnly(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.jsonl")
	if err := os.WriteFile(src, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "claude_session.jsonl")

	if _, err := copyFileStreaming(src, dst); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("copied transcript mode = %o, want 600", got)
	}
}
