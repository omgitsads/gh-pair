package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/omgitsads/gh-pair/internal/config"
	"github.com/omgitsads/gh-pair/internal/github"
)

var addCmd = &cobra.Command{
	Use:   "add <@username>",
	Short: "Add a pair by GitHub username",
	Long: `Add a GitHub user as a co-author for your commits.
The user's name and email will be fetched from GitHub.

Examples:
  gh pair add @octocat
  gh pair add octocat`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkGitRepo(); err != nil {
			return err
		}

		username := args[0]

		// Lookup user on GitHub
		pair, err := github.LookupUser(username)
		if err != nil {
			return err
		}

		// Add to config
		if err := config.AddPair(*pair); err != nil {
			return fmt.Errorf("failed to add pair: %w", err)
		}

		fmt.Printf("✓ Added: %s <%s>\n", pair.Name, pair.Email)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
