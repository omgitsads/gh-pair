package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the TUI application with the default theme.
func Run() error {
	return RunWithTheme("default")
}

// RunWithTheme starts the TUI application with the specified theme.
func RunWithTheme(themeName string) error {
	m := NewModelWithTheme(themeName)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}
