package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/omgitsads/gh-pair/internal/hook"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Install the git commit hook",
	Long: `Install the commit-msg hook in the current repository.
This hook automatically adds Co-Authored-By trailers to commits
based on your configured pairs.

If a commit-msg hook already exists, it will be backed up.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkGitRepo(); err != nil {
			return err
		}

		// Check for and remove old prepare-commit-msg hook (migration)
		if hook.HasOldHook() {
			if err := hook.RemoveOldHook(); err != nil {
				return fmt.Errorf("failed to remove old hook: %w", err)
			}
			fmt.Println("✓ Removed old prepare-commit-msg hook")
		}

		if hook.IsInstalled() {
			fmt.Println("✓ Hook already installed")
			return nil
		}

		if err := hook.Install(); err != nil {
			return fmt.Errorf("failed to install hook: %w", err)
		}

		fmt.Println("✓ Hook installed successfully")
		fmt.Println("  Use 'gh pair add @username' to add pairs")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
