package ticket

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFallbackCredentialFileAndItsDirAreOwnerOnly(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", t.TempDir())

	s := Store{Service: "fleet-test", Label: "fleet test", FileName: "test-cred.json"}
	if err := s.Save([]byte(`{"token":"x"}`)); err != nil {
		t.Fatal(err)
	}

	for path, want := range map[string]os.FileMode{
		filepath.Join(home, ".config", "fleet"):                   0o700,
		filepath.Join(home, ".config", "fleet", "test-cred.json"): 0o600,
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
