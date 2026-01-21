package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/omgitsads/gh-pair/internal/git"
	"github.com/omgitsads/gh-pair/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "gh-pair",
	Short: "Manage pair programming co-authors for git commits",
	Long: `gh-pair is a GitHub CLI extension that helps you manage 
co-authors for pair programming. It configures a git hook to 
automatically add Co-Authored-By trailers to your commits.

Run without arguments to launch the interactive TUI, or use
subcommands for quick operations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if we're in a git repo
		if !git.IsInsideWorkTree() {
			return fmt.Errorf("not a git repository")
		}

		// Launch the TUI
		return tui.Run()
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// checkGitRepo is a helper that verifies we're in a git repository.
func checkGitRepo() error {
	if !git.IsInsideWorkTree() {
		fmt.Fprintln(os.Stderr, "Error: not a git repository")
		os.Exit(1)
	}
	return nil
}
