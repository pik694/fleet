package migration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCreatesFleetConfigDirAndMarkerOwnerOnly(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("CLAUDE_CONFIG_DIR", filepath.Join(tmp, "claude"))

	prev := tmuxRunner
	tmuxRunner = &fakeTmux{sessions: []string{}}
	t.Cleanup(func() { tmuxRunner = prev })

	Run()

	newDir := filepath.Join(tmp, ".config", "fleet")
	for path, want := range map[string]os.FileMode{
		newDir: 0o700,
		filepath.Join(newDir, ".migrated-from-brizz-code"): 0o600,
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
