package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/omgitsads/gh-pair/internal/config"
)

var removeCmd = &cobra.Command{
	Use:     "remove <@username>",
	Aliases: []string{"rm"},
	Short:   "Remove a pair by GitHub username",
	Long: `Remove a GitHub user from your co-authors list.

Examples:
  gh pair remove @octocat
  gh pair rm octocat`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkGitRepo(); err != nil {
			return err
		}

		username := strings.TrimPrefix(args[0], "@")

		// Load pairs to check if exists and get display info
		pairs, err := config.LoadPairs()
		if err != nil {
			return fmt.Errorf("failed to load pairs: %w", err)
		}

		var found *config.Pair
		for _, p := range pairs.Pairs {
			if p.Username == username {
				found = &p
				break
			}
		}

		if found == nil {
			return fmt.Errorf("pair not found: @%s", username)
		}

		if err := config.RemovePair(username); err != nil {
			return fmt.Errorf("failed to remove pair: %w", err)
		}

		fmt.Printf("✓ Removed: %s <%s>\n", found.Name, found.Email)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
