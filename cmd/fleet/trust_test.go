package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brizzai/fleet/internal/workspace"
)

func repoWithFleetConfig(t *testing.T, config string) string {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, ".fleet.json"), []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	return repo
}

func TestTrustCommandRecordsTrustAndPrintsWhatItTrusted(t *testing.T) {
	repo := repoWithFleetConfig(t, `{"workspace":{"create":"make ws","destroy":"make rm"}}`)
	var out bytes.Buffer

	if err := trustRepoCommands(repo, &out); err != nil {
		t.Fatalf("trustRepoCommands: %v", err)
	}

	if !workspace.CommandsTrusted(repo) {
		t.Error("commands are still untrusted")
	}
	for _, want := range []string{"make ws", "make rm"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output %q does not show the trusted command %q", out.String(), want)
		}
	}
}

func TestTrustCommandFailsWhenTheRepoDefinesNoCommands(t *testing.T) {
	repo := repoWithFleetConfig(t, `{"pr_checks":{"ignore":["lint"]}}`)

	if err := trustRepoCommands(repo, &bytes.Buffer{}); err == nil {
		t.Fatal("trusting a repo with no shell commands reported success")
	}
}

func TestUntrustedNoticeShowsTheCommandsAndNamesTheFix(t *testing.T) {
	repo := repoWithFleetConfig(t, `{"workspace":{"create":"make ws"}}`)

	notice := untrustedCommandsNotice(repo)

	for _, want := range []string{"make ws", "fleet trust " + repo} {
		if !strings.Contains(notice, want) {
			t.Errorf("notice %q is missing %q", notice, want)
		}
	}
}
