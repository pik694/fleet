package debuglog

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func initInHome(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)
	prevLogger := Logger
	initOnce = sync.Once{}
	t.Cleanup(func() {
		Close()
		mu.Lock()
		Logger = prevLogger
		mu.Unlock()
		initOnce = sync.Once{}
	})
	Init()
}

func assertMode(t *testing.T, path string, want os.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	if got := info.Mode().Perm(); got != want {
		t.Errorf("%s mode = %o, want %o", path, got, want)
	}
}

func TestInitCreatesDebugLogAndItsDirOwnerOnly(t *testing.T) {
	home := t.TempDir()
	initInHome(t, home)

	assertMode(t, filepath.Join(home, ".config", "fleet"), 0o700)
	assertMode(t, filepath.Join(home, ".config", "fleet", "debug.log"), 0o600)
}

func TestInitTightensWorldReadableDebugLogAndDirFromAnOlderInstall(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".config", "fleet")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "debug.log"), []byte("old line\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	initInHome(t, home)

	assertMode(t, dir, 0o700)
	assertMode(t, filepath.Join(dir, "debug.log"), 0o600)
}

func TestInitTightensAnOversizedWorldReadableDebugLogItTruncates(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".config", "fleet")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	big := bytes.Repeat([]byte("0123456789abcdef\n"), maxLogSize/16)
	if err := os.WriteFile(filepath.Join(dir, "debug.log"), big, 0o644); err != nil {
		t.Fatal(err)
	}

	initInHome(t, home)

	info, err := os.Stat(filepath.Join(dir, "debug.log"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() > keepBytes {
		t.Fatalf("log was not truncated: %d bytes", info.Size())
	}
	assertMode(t, filepath.Join(dir, "debug.log"), 0o600)
}

func TestInitLeavesTheSharedTempDirModeAloneWhenHomeIsUnknown(t *testing.T) {
	tmp := t.TempDir()
	if err := os.Chmod(tmp, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMPDIR", tmp)

	initInHome(t, "")

	assertMode(t, tmp, 0o755)
	assertMode(t, filepath.Join(tmp, "fleet-debug.log"), 0o600)
}
