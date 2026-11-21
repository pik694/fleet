package stats

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenCreatesStatsDBAndItsSidecarsOwnerOnly(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fleet", "stats.db")

	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })

	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != 0o600 {
			t.Errorf("%s mode = %o, want 600", e.Name(), got)
		}
	}
}

func TestOpenTightensWorldReadableStatsDBFromAnOlderInstall(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stats.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	s.Close()
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}

	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("stats.db mode = %o, want 600", got)
	}
}
