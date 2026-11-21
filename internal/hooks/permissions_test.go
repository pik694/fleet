package hooks

import (
	"os"
	"path/filepath"
	"testing"
)

func assertPerm(t *testing.T, path string, want os.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if got := info.Mode().Perm(); got != want {
		t.Errorf("%s mode = %o, want %o", path, got, want)
	}
}

func TestWriteStatusFileCreatesHooksDirAndStatusFileOwnerOnly(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "hooks")

	if err := WriteStatusFile(dir, "abc12345-1", &StatusFile{Status: "running", UserPrompt: "secret prompt"}); err != nil {
		t.Fatal(err)
	}

	assertPerm(t, dir, 0o700)
	assertPerm(t, filepath.Join(dir, "abc12345-1.json"), 0o600)
}

func TestWriteStatusFileTightensWorldReadableHooksDirFromAnOlderInstall(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "hooks")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := WriteStatusFile(dir, "abc12345-1", &StatusFile{Status: "running"}); err != nil {
		t.Fatal(err)
	}

	assertPerm(t, dir, 0o700)
}

func TestNewHookWatcherCreatesHooksDirOwnerOnly(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	w, err := NewHookWatcher()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(w.Stop)

	assertPerm(t, filepath.Join(home, ".config", "fleet", "hooks"), 0o700)
}

func TestNewHookWatcherTightensWorldReadableHooksDirFromAnOlderInstall(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, ".config", "fleet", "hooks")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	w, err := NewHookWatcher()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(w.Stop)

	assertPerm(t, dir, 0o700)
}
