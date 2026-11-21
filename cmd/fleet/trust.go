package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/brizzai/fleet/internal/workspace"
)

const trustUsage = "Usage: fleet trust [path]"

var errNoCommandsToTrust = errors.New("this repo's .fleet.json defines no workspace shell commands; nothing to trust")

func runTrust(args []string) {
	if len(args) > 1 || (len(args) == 1 && strings.HasPrefix(args[0], "-")) {
		printTrustUsage(os.Stderr)
		os.Exit(1)
	}
	path := "."
	if len(args) == 1 {
		path = args[0]
	}
	repoPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot resolve %s: %v\n", path, err)
		os.Exit(1)
	}
	if err := trustRepoCommands(repoPath, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printTrustUsage(w io.Writer) {
	fmt.Fprintln(w, trustUsage)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Allows fleet to run the workspace shell commands (list/create/destroy) that")
	fmt.Fprintln(w, "the repo's .fleet.json defines. Trust covers exactly the commands printed;")
	fmt.Fprintln(w, "if they change, fleet refuses to run them until you trust them again.")
}

func trustRepoCommands(repoPath string, out io.Writer) error {
	commands := workspace.RepoCommands(repoPath)
	if len(commands.Lines()) == 0 {
		return errNoCommandsToTrust
	}
	if err := workspace.TrustCommands(repoPath); err != nil {
		return fmt.Errorf("could not record trust: %w", err)
	}
	fmt.Fprintf(out, "Trusted the .fleet.json workspace commands in %s:\n%s\n", repoPath, indentedCommands(commands))
	return nil
}

func untrustedCommandsNotice(repoPath string) string {
	return fmt.Sprintf(
		"This repo's .fleet.json defines workspace shell commands that you have not trusted:\n%s\nReview them, then run: fleet trust %s",
		indentedCommands(workspace.RepoCommands(repoPath)), repoPath)
}

func indentedCommands(commands workspace.ShellConfig) string {
	return "  " + strings.Join(commands.Lines(), "\n  ")
}
