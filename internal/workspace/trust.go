package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func CommandsTrusted(repoPath string) bool {
	commands := loadMergedRepoConfig(repoPath).Workspace
	if !commands.any() {
		return true
	}
	return loadTrusted()[mainCheckout(repoPath)] == commands.fingerprint()
}

func TrustCommands(repoPath string) error {
	trusted := loadTrusted()
	trusted[mainCheckout(repoPath)] = loadMergedRepoConfig(repoPath).Workspace.fingerprint()
	return saveTrusted(trusted)
}

func (c ShellConfig) any() bool {
	return c.List != "" || c.Create != "" || c.Destroy != ""
}

func (c ShellConfig) fingerprint() string {
	sum := sha256.Sum256([]byte(c.List + "\x00" + c.Create + "\x00" + c.Destroy))
	return hex.EncodeToString(sum[:])
}

func trustFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "fleet", "trusted_commands.json")
}

func loadTrusted() map[string]string {
	trusted := map[string]string{}
	data, err := os.ReadFile(trustFilePath())
	if err != nil {
		return trusted
	}
	_ = json.Unmarshal(data, &trusted)
	return trusted
}

func saveTrusted(trusted map[string]string) error {
	data, err := json.MarshalIndent(trusted, "", "  ")
	if err != nil {
		return err
	}
	path := trustFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func mainCheckout(repoPath string) string {
	pointer, err := os.ReadFile(filepath.Join(repoPath, ".git"))
	if err != nil {
		return repoPath
	}
	gitDir, ok := strings.CutPrefix(strings.TrimSpace(string(pointer)), "gitdir:")
	if !ok {
		return repoPath
	}
	gitDir = resolveAgainst(repoPath, strings.TrimSpace(gitDir))
	commonDir, err := os.ReadFile(filepath.Join(gitDir, "commondir"))
	if err != nil {
		return repoPath
	}
	return filepath.Dir(resolveAgainst(gitDir, strings.TrimSpace(string(commonDir))))
}

func resolveAgainst(base, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(base, path)
}
