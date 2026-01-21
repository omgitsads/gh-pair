package tui

import (
	"fmt"
	"strings"
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
	b.WriteString(m.styles.Title.Render("🤝 gh-pair"))
	b.WriteString("\n")

	// Hook status
	if !m.hookInstalled {
		b.WriteString(m.styles.Error.Render("⚠ Hook not installed"))
		b.WriteString(m.styles.Dim.Render(" - press 'i' to install"))
		b.WriteString("\n\n")
	}

	// Error display
	if m.err != nil {
		b.WriteString(m.styles.Error.Render("Error: " + m.err.Error()))
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
		b.WriteString(m.styles.Subtitle.Render("No pairs configured"))
		b.WriteString("\n")
		b.WriteString(m.styles.Dim.Render("Press 'a' to search users or 't' to browse teams"))
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

	b.WriteString(m.styles.Title.Render("🔍 Add Pair"))
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
		b.WriteString(m.styles.Error.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	// Results label
	if len(m.searchResults) > 0 {
		b.WriteString(m.styles.Dim.Render("Search Results:"))
	} else if len(m.searchList.Items()) > 0 {
		b.WriteString(m.styles.Dim.Render("Recent / Collaborators:"))
	}
	b.WriteString("\n")

	// Search results list
	b.WriteString(m.searchList.View())

	// Help footer
	b.WriteString("\n")
	b.WriteString(m.styles.Dim.Render("Enter: add • Tab: switch focus • Esc: cancel"))

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
	content.WriteString(m.styles.Title.Render("Keyboard Shortcuts"))
	content.WriteString("\n\n")

	for _, h := range help {
		content.WriteString(fmt.Sprintf("%s  %s\n",
			m.styles.HelpKey.Render(fmt.Sprintf("%-12s", h.key)),
			m.styles.HelpDesc.Render(h.desc)))
	}

	content.WriteString("\n")
	content.WriteString(m.styles.Dim.Render("Press Esc or ? to close"))

	return m.styles.Box.Render(content.String())
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
			m.styles.HelpKey.Render(k.key),
			m.styles.Dim.Render(k.desc)))
	}

	return strings.Join(parts, m.styles.Dim.Render(" • "))
}

func (m Model) teamsView() string {
	var b strings.Builder

	b.WriteString(m.styles.Title.Render("👥 Your Teams"))
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
		b.WriteString(m.styles.Error.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	if len(m.filteredTeams) == 0 && len(m.teams) > 0 {
		b.WriteString(m.styles.Dim.Render("No teams match your filter"))
		b.WriteString("\n")
	} else if len(m.teams) == 0 {
		b.WriteString(m.styles.Subtitle.Render("No teams found"))
		b.WriteString("\n")
		b.WriteString(m.styles.Dim.Render("You may not be a member of any GitHub teams"))
		b.WriteString("\n")
	} else {
		b.WriteString(m.teamList.View())
	}

	b.WriteString("\n")
	b.WriteString(m.styles.Dim.Render("Enter: select team • Tab: filter • Esc: back"))

	return b.String()
}

func (m Model) teamMembersView() string {
	var b strings.Builder

	teamName := ""
	if m.selectedTeam != nil {
		teamName = m.selectedTeam.Name
	}
	b.WriteString(m.styles.Title.Render("👥 " + teamName + " Members"))
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
		b.WriteString(m.styles.Error.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	if len(m.filteredTeamMembers) == 0 && len(m.teamMembers) > 0 {
		b.WriteString(m.styles.Dim.Render("No members match your filter"))
		b.WriteString("\n")
	} else if len(m.teamMembers) == 0 {
		b.WriteString(m.styles.Dim.Render("No members found"))
		b.WriteString("\n")
	} else {
		b.WriteString(m.searchList.View())
	}

	b.WriteString("\n")
	b.WriteString(m.styles.Dim.Render("Enter: add • Tab: switch focus • Esc: back to teams"))

	return b.String()
}
