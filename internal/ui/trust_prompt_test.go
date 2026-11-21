package ui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brizzai/fleet/internal/config"
	"github.com/brizzai/fleet/internal/session"
	"github.com/brizzai/fleet/internal/workspace"
)

func homeWithUntrustedRepo(t *testing.T, fleetConfig string) (*Home, string) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	storage, err := session.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { storage.Close() })
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, ".fleet.json"), []byte(fleetConfig), 0o644); err != nil {
		t.Fatal(err)
	}
	home := NewHome(storage, &config.Config{TickIntervalSec: 2}, "test")
	home.width, home.height = 120, 40
	return home, repo
}

func TestWorktreeListFetchReportsUntrustedCommandsBeforeRunningAnything(t *testing.T) {
	home, repo := homeWithUntrustedRepo(t, `{"workspace":{"create":"touch ran"}}`)

	msg := home.fetchWorkspaceListForRepo(repo)().(workspaceListMsg)

	if !errors.Is(msg.err, workspace.ErrUntrustedCommands) {
		t.Fatalf("err = %v, want ErrUntrustedCommands", msg.err)
	}
	if msg.repoPath != repo {
		t.Fatalf("repoPath = %q, want %q", msg.repoPath, repo)
	}
}

func TestUntrustedWorktreeListOpensATrustPromptShowingTheCommands(t *testing.T) {
	home, repo := homeWithUntrustedRepo(t, `{"workspace":{"create":"make ws"}}`)

	home.Update(workspaceListMsg{repoPath: repo, err: workspace.ErrUntrustedCommands})

	if !home.confirmDialog.IsVisible() {
		t.Fatal("no trust prompt was shown")
	}
	if view := home.confirmDialog.View(); !strings.Contains(view, "make ws") {
		t.Fatalf("trust prompt does not show the command:\n%s", view)
	}
	if home.worktreeDialog.IsVisible() {
		t.Fatal("the worktree dialog stayed open behind the trust prompt")
	}
}

func TestConfirmingTheTrustPromptTrustsTheCommands(t *testing.T) {
	home, repo := homeWithUntrustedRepo(t, `{"workspace":{"create":"make ws"}}`)
	home.Update(workspaceListMsg{repoPath: repo, err: workspace.ErrUntrustedCommands})

	home.Update(home.confirmDialog.onYes())

	if !workspace.CommandsTrusted(repo) {
		t.Fatal("confirming the prompt left the commands untrusted")
	}
}

func TestRemovingAWorktreeWithUntrustedCommandsAsksForTrustInstead(t *testing.T) {
	home, repo := homeWithUntrustedRepo(t, `{"workspace":{"destroy":"make rm"}}`)
	home.failedWorktreeRemovals = map[string]bool{repo: true}

	home.confirmDeleteHeader(SidebarItem{IsRepoHeader: true, RepoPath: repo})

	view := home.confirmDialog.View()
	if !strings.Contains(view, "make rm") {
		t.Fatalf("expected a trust prompt showing the destroy command, got:\n%s", view)
	}
	if strings.Contains(view, "Remove Worktree?") {
		t.Fatalf("the removal prompt was offered for untrusted commands:\n%s", view)
	}
}
