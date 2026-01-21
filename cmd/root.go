package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/omgitsads/gh-pair/internal/git"
	"github.com/omgitsads/gh-pair/internal/theme"
	"github.com/omgitsads/gh-pair/internal/tui"
)

var themeName string
var themeFlag bool // tracks if --theme was explicitly set

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

		// Launch the TUI with theme
		return tui.RunWithTheme(getThemeName())
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&themeName, "theme", "", "Color theme (default, dracula, nord, solarized-dark, solarized-light, catppuccin)")
}

// getThemeName returns the theme name from flag or config.
func getThemeName() string {
	if themeName != "" {
		return themeName
	}
	return theme.GetConfiguredTheme()
}

// checkGitRepo is a helper that verifies we're in a git repository.
func checkGitRepo() error {
	if !git.IsInsideWorkTree() {
		fmt.Fprintln(os.Stderr, "Error: not a git repository")
		os.Exit(1)
	}
	return nil
}
