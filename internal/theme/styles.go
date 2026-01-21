package theme

import "github.com/charmbracelet/lipgloss"

// Styles contains all lipgloss styles derived from a theme.
type Styles struct {
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Error       lipgloss.Style
	Success     lipgloss.Style
	Warning     lipgloss.Style
	Dim         lipgloss.Style
	HelpKey     lipgloss.Style
	HelpDesc    lipgloss.Style
	Box         lipgloss.Style
	ListTitle   lipgloss.Style
	Spinner     lipgloss.Style
}

// NewStyles creates a Styles instance from a theme.
func NewStyles(t Theme) Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.Colors.Primary)).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.TextDim)).
			MarginBottom(1),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Error)).
			Bold(true),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Success)),

		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Warning)),

		Dim: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.TextDim)),

		HelpKey: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Secondary)).
			Bold(true),

		HelpDesc: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.TextDim)),

		Box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.Colors.Border)).
			Padding(1, 2),

		ListTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(t.Colors.Primary)).
			Padding(0, 1),

		Spinner: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Secondary)),
	}
}

// DefaultStyles returns styles using the default theme.
func DefaultStyles() Styles {
	return NewStyles(DefaultTheme)
}
