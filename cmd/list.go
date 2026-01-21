package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/omgitsads/gh-pair/internal/config"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List current pairs",
	Long:    `Display all currently configured co-authors.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkGitRepo(); err != nil {
			return err
		}

		pairs, err := config.LoadPairs()
		if err != nil {
			return fmt.Errorf("failed to load pairs: %w", err)
		}

		if len(pairs.Pairs) == 0 {
			fmt.Println("No pairs configured")
			fmt.Println("Use 'gh pair add @username' to add pairs")
			return nil
		}

		fmt.Println("Current pairs:")
		for _, p := range pairs.Pairs {
			fmt.Printf("  @%-20s %s <%s>\n", p.Username, p.Name, p.Email)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
