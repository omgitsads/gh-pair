package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginBottom(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(1, 2)
)

// View renders the TUI.
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	switch m.view {
	case ViewHelp:
		return m.helpView()
	case ViewSearch:
		return m.searchView()
	default:
		return m.mainView()
	}
}

func (m Model) mainView() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("🤝 gh-pair"))
	b.WriteString("\n")

	// Hook status
	if !m.hookInstalled {
		b.WriteString(errorStyle.Render("⚠ Hook not installed"))
		b.WriteString(dimStyle.Render(" - press 'i' to install"))
		b.WriteString("\n\n")
	}

	// Error display
	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	// Loading indicator
	if m.loading {
		b.WriteString(m.spinner.View())
		b.WriteString(" Loading...\n")
		return b.String()
	}

	// Pair list or empty state
	if len(m.pairs) == 0 {
		b.WriteString(subtitleStyle.Render("No pairs configured"))
		b.WriteString("\n")
		b.WriteString(dimStyle.Render("Press 'a' to add a pair"))
		b.WriteString("\n")
	} else {
		b.WriteString(m.pairList.View())
	}

	// Help footer
	b.WriteString("\n")
	b.WriteString(m.helpFooter())

	return b.String()
}

func (m Model) searchView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🔍 Add Pair"))
	b.WriteString("\n\n")

	// Search input
	b.WriteString(m.searchInput.View())
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(m.spinner.View())
		b.WriteString(" Searching...\n")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	// Results label
	if len(m.searchResults) > 0 {
		b.WriteString(dimStyle.Render("Search Results:"))
	} else if len(m.searchList.Items()) > 0 {
		b.WriteString(dimStyle.Render("Recent / Collaborators:"))
	}
	b.WriteString("\n")

	// Search results list
	b.WriteString(m.searchList.View())

	// Help footer
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("Enter: add • Tab: switch focus • Esc: cancel"))

	return b.String()
}

func (m Model) helpView() string {
	help := []struct {
		key  string
		desc string
	}{
		{"a, /", "Add a new pair"},
		{"d, Delete", "Remove selected pair"},
		{"c", "Clear all pairs"},
		{"i", "Install git hook"},
		{"↑/↓", "Navigate list"},
		{"Enter", "Select / Confirm"},
		{"Esc", "Cancel / Back"},
		{"?", "Toggle help"},
		{"q", "Quit"},
	}

	var content strings.Builder
	content.WriteString(titleStyle.Render("Keyboard Shortcuts"))
	content.WriteString("\n\n")

	for _, h := range help {
		content.WriteString(fmt.Sprintf("%s  %s\n",
			helpKeyStyle.Render(fmt.Sprintf("%-12s", h.key)),
			helpDescStyle.Render(h.desc)))
	}

	content.WriteString("\n")
	content.WriteString(dimStyle.Render("Press Esc or ? to close"))

	return boxStyle.Render(content.String())
}

func (m Model) helpFooter() string {
	keys := []struct {
		key  string
		desc string
	}{
		{"a", "add"},
		{"d", "remove"},
		{"c", "clear"},
		{"?", "help"},
		{"q", "quit"},
	}

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s %s",
			helpKeyStyle.Render(k.key),
			dimStyle.Render(k.desc)))
	}

	return strings.Join(parts, dimStyle.Render(" • "))
}
