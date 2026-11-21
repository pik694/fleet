package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

const everyCommandLeavesAMarker = `{"workspace":{"list":"touch ran; echo '[]'","create":"touch ran","destroy":"touch ran"}}`

func runEveryShellOperation(provider Provider, repo string) []error {
	_, listErr := provider.List(repo)
	_, createErr := provider.Create(repo, "ws", "branch", "")
	destroyErr := provider.Destroy(repo, "ws")
	return []error{listErr, createErr, destroyErr}
}

func TestShellProviderRefusesEveryOperationWhileUntrusted(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	repo := t.TempDir()
	writeFile(t, repo, ".fleet.json", everyCommandLeavesAMarker)

	for _, err := range runEveryShellOperation(ResolveProvider(repo), repo) {
		if !errors.Is(err, ErrUntrustedCommands) {
			t.Errorf("err = %v, want ErrUntrustedCommands", err)
		}
	}
	if _, err := os.Stat(filepath.Join(repo, "ran")); err == nil {
		t.Fatal("an untrusted command was executed")
	}
}

func TestShellProviderRunsEveryOperationOnceTrusted(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	repo := t.TempDir()
	writeFile(t, repo, ".fleet.json", everyCommandLeavesAMarker)
	if err := TrustCommands(repo); err != nil {
		t.Fatalf("TrustCommands: %v", err)
	}

	for _, err := range runEveryShellOperation(ResolveProvider(repo), repo) {
		if err != nil {
			t.Errorf("trusted operation failed: %v", err)
		}
	}
}

func TestShellProviderRefusesCommandsThatDifferFromTheTrustedOnes(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	repo := t.TempDir()
	writeFile(t, repo, ".fleet.json", `{"workspace":{"destroy":"true"}}`)
	if err := TrustCommands(repo); err != nil {
		t.Fatalf("TrustCommands: %v", err)
	}
	swapped := &ShellProvider{DestroyCmd: "touch ran"}

	if err := swapped.Destroy(repo, "ws"); !errors.Is(err, ErrUntrustedCommands) {
		t.Fatalf("err = %v, want ErrUntrustedCommands", err)
	}
}
