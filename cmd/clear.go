package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/omgitsads/gh-pair/internal/config"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all pairs",
	Long:  `Clear all co-authors from the current repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkGitRepo(); err != nil {
			return err
		}

		if err := config.ClearPairs(); err != nil {
			return fmt.Errorf("failed to clear pairs: %w", err)
		}

		fmt.Println("✓ All pairs cleared")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
