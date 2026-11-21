package ui

import (
	"fmt"
	"path/filepath"

	tea "charm.land/bubbletea/v2"

	"github.com/brizzai/fleet/internal/workspace"
)

type trustCommandsMsg struct {
	repoPath string
}

func (h *Home) askToTrustCommands(repoPath string) {
	details := append(workspace.RepoCommands(repoPath).Lines(),
		"Defined in .fleet.json — fleet runs them in your shell",
		"Asked again whenever they change")
	h.confirmDialog.ShowWarning("Trust this repo's workspace commands?", filepath.Base(repoPath), details, func() tea.Msg {
		return trustCommandsMsg{repoPath: repoPath}
	})
}

func (h *Home) trustCommands(repoPath string) {
	if err := workspace.TrustCommands(repoPath); err != nil {
		h.setError(fmt.Errorf("trust commands: %w", err))
		return
	}
	h.setInfo("Commands trusted — repeat the action")
}
