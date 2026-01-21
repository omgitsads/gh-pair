package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/omgitsads/gh-pair/internal/config"
	"github.com/omgitsads/gh-pair/internal/github"
	"github.com/omgitsads/gh-pair/internal/hook"
)

// View represents the current view in the TUI.
type View int

const (
	ViewMain View = iota
	ViewSearch
	ViewHelp
)

// Model is the main TUI model.
type Model struct {
	view          View
	pairs         []config.Pair
	recentPairs   []config.Pair
	collaborators []config.Pair
	searchResults []config.Pair

	pairList      list.Model
	searchInput   textinput.Model
	searchList    list.Model
	spinner       spinner.Model
	loading       bool
	hookInstalled bool
	err           error

	width  int
	height int
}

// pairItem implements list.Item for pairs.
type pairItem struct {
	pair config.Pair
}

func (i pairItem) Title() string       { return "@" + i.pair.Username }
func (i pairItem) Description() string { return i.pair.Name + " <" + i.pair.Email + ">" }
func (i pairItem) FilterValue() string { return i.pair.Username + " " + i.pair.Name }

// Messages
type (
	pairsLoadedMsg struct {
		pairs  []config.Pair
		recent []config.Pair
	}
	collaboratorsLoadedMsg struct {
		collaborators []config.Pair
	}
	searchResultsMsg struct {
		results []config.Pair
	}
	userLookedUpMsg struct {
		pair *config.Pair
		err  error
	}
	errMsg struct {
		err error
	}
)

// NewModel creates a new TUI model.
func NewModel() Model {
	// Set up spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Set up search input
	ti := textinput.New()
	ti.Placeholder = "Search GitHub users..."
	ti.CharLimit = 50
	ti.Width = 40

	// Set up pair list
	delegate := list.NewDefaultDelegate()
	pairList := list.New([]list.Item{}, delegate, 0, 0)
	pairList.Title = "Current Pairs"
	pairList.SetShowStatusBar(false)
	pairList.SetFilteringEnabled(false)
	pairList.Styles.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Padding(0, 1)

	// Set up search results list
	searchList := list.New([]list.Item{}, delegate, 0, 0)
	searchList.Title = "Search Results"
	searchList.SetShowStatusBar(false)
	searchList.SetFilteringEnabled(false)

	return Model{
		view:        ViewMain,
		pairList:    pairList,
		searchInput: ti,
		searchList:  searchList,
		spinner:     s,
		loading:     true,
	}
}

// Init initializes the TUI.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadPairs,
		loadCollaborators,
	)
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.pairList.SetSize(msg.Width-4, msg.Height-8)
		m.searchList.SetSize(msg.Width-4, msg.Height-12)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case pairsLoadedMsg:
		m.pairs = msg.pairs
		m.recentPairs = msg.recent
		m.loading = false
		m.hookInstalled = hook.IsInstalled()
		m.updatePairList()
		return m, nil

	case collaboratorsLoadedMsg:
		m.collaborators = msg.collaborators
		return m, nil

	case searchResultsMsg:
		m.searchResults = msg.results
		m.loading = false
		m.updateSearchList()
		return m, nil

	case userLookedUpMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		if msg.pair != nil {
			if err := config.AddPair(*msg.pair); err != nil {
				m.err = err
				return m, nil
			}
			m.view = ViewMain
			return m, loadPairs
		}
		return m, nil

	case errMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	}

	// Update sub-models
	switch m.view {
	case ViewMain:
		var cmd tea.Cmd
		m.pairList, cmd = m.pairList.Update(msg)
		cmds = append(cmds, cmd)
	case ViewSearch:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
		m.searchList, cmd = m.searchList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys
	switch msg.String() {
	case "ctrl+c", "q":
		if m.view == ViewMain {
			return m, tea.Quit
		}
		m.view = ViewMain
		m.searchInput.SetValue("")
		m.searchResults = nil
		return m, nil

	case "esc":
		if m.view != ViewMain {
			m.view = ViewMain
			m.searchInput.SetValue("")
			m.searchResults = nil
			return m, nil
		}
		return m, tea.Quit

	case "?":
		if m.view == ViewMain {
			m.view = ViewHelp
			return m, nil
		}
		if m.view == ViewHelp {
			m.view = ViewMain
			return m, nil
		}
	}

	// View-specific keys
	switch m.view {
	case ViewMain:
		return m.handleMainKeys(msg)
	case ViewSearch:
		return m.handleSearchKeys(msg)
	case ViewHelp:
		if msg.String() == "enter" || msg.String() == "esc" || msg.String() == "?" {
			m.view = ViewMain
			return m, nil
		}
	}

	return m, nil
}

func (m Model) handleMainKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "a", "/":
		m.view = ViewSearch
		m.searchInput.Focus()
		m.updateSearchList() // Show recent/collaborators initially
		return m, nil

	case "d", "backspace", "delete":
		if item, ok := m.pairList.SelectedItem().(pairItem); ok {
			if err := config.RemovePair(item.pair.Username); err != nil {
				m.err = err
				return m, nil
			}
			return m, loadPairs
		}

	case "c":
		if err := config.ClearPairs(); err != nil {
			m.err = err
			return m, nil
		}
		return m, loadPairs

	case "i":
		if !m.hookInstalled {
			if err := hook.Install(); err != nil {
				m.err = err
				return m, nil
			}
			m.hookInstalled = true
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.pairList, cmd = m.pairList.Update(msg)
	return m, cmd
}

func (m Model) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.searchInput.Focused() {
			query := strings.TrimSpace(m.searchInput.Value())
			if query != "" {
				// If it looks like a username, try direct lookup
				if strings.HasPrefix(query, "@") || !strings.Contains(query, " ") {
					m.loading = true
					return m, lookupUser(query)
				}
				// Otherwise search
				m.loading = true
				return m, searchUsers(query)
			}
		} else {
			// Select from search list
			if item, ok := m.searchList.SelectedItem().(pairItem); ok {
				if err := config.AddPair(item.pair); err != nil {
					m.err = err
					return m, nil
				}
				m.view = ViewMain
				m.searchInput.SetValue("")
				m.searchResults = nil
				return m, loadPairs
			}
		}

	case "tab":
		if m.searchInput.Focused() {
			m.searchInput.Blur()
		} else {
			m.searchInput.Focus()
		}
		return m, nil

	case "up", "down":
		if !m.searchInput.Focused() {
			var cmd tea.Cmd
			m.searchList, cmd = m.searchList.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m *Model) updatePairList() {
	items := make([]list.Item, len(m.pairs))
	for i, p := range m.pairs {
		items[i] = pairItem{pair: p}
	}
	m.pairList.SetItems(items)
}

func (m *Model) updateSearchList() {
	var items []list.Item

	if len(m.searchResults) > 0 {
		for _, p := range m.searchResults {
			items = append(items, pairItem{pair: p})
		}
	} else {
		// Show recent pairs and collaborators
		seen := make(map[string]bool)
		for _, p := range m.pairs {
			seen[p.Username] = true
		}

		for _, p := range m.recentPairs {
			if !seen[p.Username] {
				items = append(items, pairItem{pair: p})
				seen[p.Username] = true
			}
		}

		for _, p := range m.collaborators {
			if !seen[p.Username] {
				items = append(items, pairItem{pair: p})
				seen[p.Username] = true
			}
		}
	}

	m.searchList.SetItems(items)
}

// Commands
func loadPairs() tea.Msg {
	pairs, err := config.LoadPairs()
	if err != nil {
		return errMsg{err: err}
	}
	recent, _ := config.LoadRecent()
	return pairsLoadedMsg{pairs: pairs.Pairs, recent: recent.Recent}
}

func loadCollaborators() tea.Msg {
	collaborators, _ := github.GetRepoCollaborators()
	return collaboratorsLoadedMsg{collaborators: collaborators}
}

func searchUsers(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := github.SearchUsers(query)
		if err != nil {
			return errMsg{err: err}
		}
		return searchResultsMsg{results: results}
	}
}

func lookupUser(username string) tea.Cmd {
	return func() tea.Msg {
		pair, err := github.LookupUser(username)
		return userLookedUpMsg{pair: pair, err: err}
	}
}
