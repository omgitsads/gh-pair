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
	case ViewTeams:
		return m.teamsView()
	case ViewTeamMembers:
		return m.teamMembersView()
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
		b.WriteString(dimStyle.Render("Press 'a' to search users or 't' to browse teams"))
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
		{"a, /", "Search GitHub users"},
		{"t", "Browse your teams"},
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
		{"a", "search"},
		{"t", "teams"},
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

func (m Model) teamsView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("👥 Your Teams"))
	b.WriteString("\n\n")

	// Filter input
	b.WriteString(m.searchInput.View())
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(m.spinner.View())
		b.WriteString(" Loading teams...\n")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	if len(m.filteredTeams) == 0 && len(m.teams) > 0 {
		b.WriteString(dimStyle.Render("No teams match your filter"))
		b.WriteString("\n")
	} else if len(m.teams) == 0 {
		b.WriteString(subtitleStyle.Render("No teams found"))
		b.WriteString("\n")
		b.WriteString(dimStyle.Render("You may not be a member of any GitHub teams"))
		b.WriteString("\n")
	} else {
		b.WriteString(m.teamList.View())
	}

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("Enter: select team • Tab: filter • Esc: back"))

	return b.String()
}

func (m Model) teamMembersView() string {
	var b strings.Builder

	teamName := ""
	if m.selectedTeam != nil {
		teamName = m.selectedTeam.Name
	}
	b.WriteString(titleStyle.Render("👥 " + teamName + " Members"))
	b.WriteString("\n\n")

	// Filter input
	b.WriteString(m.searchInput.View())
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(m.spinner.View())
		b.WriteString(" Loading members...\n")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	if len(m.filteredTeamMembers) == 0 && len(m.teamMembers) > 0 {
		b.WriteString(dimStyle.Render("No members match your filter"))
		b.WriteString("\n")
	} else if len(m.teamMembers) == 0 {
		b.WriteString(dimStyle.Render("No members found"))
		b.WriteString("\n")
	} else {
		b.WriteString(m.searchList.View())
	}

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("Enter: add • Tab: switch focus • Esc: back to teams"))

	return b.String()
}
