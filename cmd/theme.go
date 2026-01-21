package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/omgitsads/gh-pair/internal/theme"
)

var themeCmd = &cobra.Command{
	Use:   "theme [name]",
	Short: "List available themes or preview a theme",
	Long: `List all available themes (preset and custom) or preview a specific theme.

Examples:
  gh pair theme              # List all available themes
  gh pair theme dracula      # Preview the dracula theme
  gh pair theme set dracula  # Set dracula as your default theme`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			listThemes()
		} else {
			previewTheme(args[0])
		}
	},
}

var themeSetCmd = &cobra.Command{
	Use:   "set <theme>",
	Short: "Set the default theme",
	Long:  `Set the default theme to use when launching gh-pair.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		// Verify theme exists
		t := theme.GetTheme(name)
		if t.Name == "" && name != "default" {
			// Theme not found, check if it's a valid preset or custom
			found := false
			for _, preset := range theme.PresetNames() {
				if preset == name {
					found = true
					break
				}
			}
			for _, custom := range theme.ListCustomThemes() {
				if custom == name {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("theme %q not found", name)
			}
		}

		if err := theme.SetConfiguredTheme(name); err != nil {
			return fmt.Errorf("failed to save theme: %w", err)
		}

		styles := theme.NewStyles(theme.GetTheme(name))
		fmt.Printf("%s Theme set to %s\n", 
			styles.Success.Render("✓"),
			styles.Title.Render(name))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(themeCmd)
	themeCmd.AddCommand(themeSetCmd)
}

func listThemes() {
	currentTheme := theme.GetConfiguredTheme()
	
	fmt.Println("Available themes:")
	fmt.Println()

	// Preset themes
	fmt.Println("Preset themes:")
	for _, name := range theme.PresetNames() {
		t := theme.GetTheme(name)
		preview := renderColorPreview(t)
		marker := "  "
		if name == currentTheme {
			marker = "► "
		}
		fmt.Printf("%s%-18s %s\n", marker, name, preview)
	}

	// Custom themes
	customThemes := theme.ListCustomThemes()
	if len(customThemes) > 0 {
		fmt.Println()
		fmt.Println("Custom themes (~/.config/gh-pair/themes/):")
		for _, name := range customThemes {
			t := theme.GetTheme(name)
			preview := renderColorPreview(t)
			marker := "  "
			if name == currentTheme {
				marker = "► "
			}
			fmt.Printf("%s%-18s %s\n", marker, name, preview)
		}
	}

	fmt.Println()
	fmt.Println("Use: gh pair theme set <name>")
}

func previewTheme(name string) {
	t := theme.GetTheme(name)
	styles := theme.NewStyles(t)

	fmt.Println()
	fmt.Println(styles.Title.Render("Theme: " + t.Name))
	fmt.Println()

	// Color swatches
	fmt.Printf("  Primary:   %s\n", renderSwatch(t.Colors.Primary, "████"))
	fmt.Printf("  Secondary: %s\n", renderSwatch(t.Colors.Secondary, "████"))
	fmt.Printf("  Success:   %s\n", renderSwatch(t.Colors.Success, "████"))
	fmt.Printf("  Error:     %s\n", renderSwatch(t.Colors.Error, "████"))
	fmt.Printf("  Warning:   %s\n", renderSwatch(t.Colors.Warning, "████"))
	fmt.Printf("  Border:    %s\n", renderSwatch(t.Colors.Border, "████"))
	fmt.Printf("  Accent:    %s\n", renderSwatch(t.Colors.Accent, "████"))
	fmt.Printf("  Text:      %s\n", renderSwatch(t.Colors.Text, "████"))
	fmt.Printf("  TextDim:   %s\n", renderSwatch(t.Colors.TextDim, "████"))
	fmt.Println()

	// Sample UI elements
	fmt.Println(styles.Title.Render("🤝 Sample Title"))
	fmt.Println(styles.Subtitle.Render("Subtitle text"))
	fmt.Println(styles.Success.Render("✓ Success message"))
	fmt.Println(styles.Error.Render("✗ Error message"))
	fmt.Println(styles.Warning.Render("⚠ Warning message"))
	fmt.Println(styles.Dim.Render("Dimmed hint text"))
	fmt.Printf("%s %s\n", styles.HelpKey.Render("a"), styles.HelpDesc.Render("help key"))
	fmt.Println()
}

func renderColorPreview(t theme.Theme) string {
	return fmt.Sprintf("%s%s%s%s%s",
		renderSwatch(t.Colors.Primary, "█"),
		renderSwatch(t.Colors.Secondary, "█"),
		renderSwatch(t.Colors.Success, "█"),
		renderSwatch(t.Colors.Error, "█"),
		renderSwatch(t.Colors.Accent, "█"))
}

func renderSwatch(color string, text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(text)
}
