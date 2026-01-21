package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/omgitsads/gh-pair/internal/config"
	"github.com/omgitsads/gh-pair/internal/github"
	"github.com/omgitsads/gh-pair/internal/hook"
)

const debounceDelay = 300 * time.Millisecond

// View represents the current view in the TUI.
type View int

const (
	ViewMain View = iota
	ViewSearch
	ViewTeams
	ViewTeamMembers
	ViewHelp
)

// SearchTab represents which tab is active in search view.
type SearchTab int

const (
	TabUsers SearchTab = iota
	TabTeams
)

// Model is the main TUI model.
type Model struct {
	view          View
	pairs         []config.Pair
	recentPairs   []config.Pair
	collaborators []config.Pair
	searchResults []config.Pair

	// Team-related state
	teams           []github.Team
	filteredTeams   []github.Team
	selectedTeam    *github.Team
	teamMembers     []config.Pair
	filteredTeamMembers []config.Pair
	searchTab       SearchTab

	pairList      list.Model
	searchInput   textinput.Model
	searchList    list.Model
	teamList      list.Model
	spinner       spinner.Model
	loading       bool
	hookInstalled bool
	err           error

	// Debounce state for autocomplete
	lastQuery     string
	debounceTimer int // incremented each time we schedule a debounce

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

// teamItem implements list.Item for teams.
type teamItem struct {
	team github.Team
}

func (i teamItem) Title() string       { return i.team.Name }
func (i teamItem) Description() string { return i.team.Org + "/" + i.team.Slug }
func (i teamItem) FilterValue() string { return i.team.Name + " " + i.team.Slug }

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
		query   string // track which query this result is for
	}
	userLookedUpMsg struct {
		pair *config.Pair
		err  error
	}
	errMsg struct {
		err error
	}
	// debounceTickMsg is sent after the debounce delay
	debounceTickMsg struct {
		query   string
		timerID int
	}
	teamsLoadedMsg struct {
		teams []github.Team
	}
	teamMembersLoadedMsg struct {
		members []config.Pair
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

	// Set up team list
	teamList := list.New([]list.Item{}, delegate, 0, 0)
	teamList.Title = "Your Teams"
	teamList.SetShowStatusBar(false)
	teamList.SetFilteringEnabled(false)

	return Model{
		view:        ViewMain,
		pairList:    pairList,
		searchInput: ti,
		searchList:  searchList,
		teamList:    teamList,
		spinner:     s,
		loading:     true,
		searchTab:   TabUsers,
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
		m.teamList.SetSize(msg.Width-4, msg.Height-12)
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
		// Only update if this result matches the current query
		if msg.query == m.lastQuery {
			m.searchResults = msg.results
			m.loading = false
			m.updateSearchList()
		}
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

	case teamsLoadedMsg:
		m.teams = msg.teams
		m.filteredTeams = msg.teams
		m.loading = false
		m.updateTeamList()
		return m, nil

	case teamMembersLoadedMsg:
		m.teamMembers = msg.members
		m.filteredTeamMembers = msg.members
		m.loading = false
		m.updateSearchList()
		return m, nil

	case debounceTickMsg:
		// Only trigger search if this is the latest timer and query matches
		if msg.timerID == m.debounceTimer && msg.query == m.searchInput.Value() {
			query := strings.TrimSpace(msg.query)
			if len(query) >= 2 {
				m.loading = true
				m.lastQuery = query
				return m, searchUsers(query)
			}
		}
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
		oldValue := m.searchInput.Value()
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)

		// Check if input changed - schedule debounced search
		newValue := m.searchInput.Value()
		if newValue != oldValue && m.searchInput.Focused() {
			m.debounceTimer++
			cmds = append(cmds, scheduleDebounce(newValue, m.debounceTimer))
		}

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
		m.selectedTeam = nil
		m.teamMembers = nil
		return m, nil

	case "esc":
		if m.view == ViewTeamMembers {
			m.view = ViewTeams
			m.selectedTeam = nil
			m.teamMembers = nil
			m.searchInput.SetValue("")
			return m, nil
		}
		if m.view != ViewMain {
			m.view = ViewMain
			m.searchInput.SetValue("")
			m.searchResults = nil
			m.selectedTeam = nil
			m.teamMembers = nil
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
	case ViewTeams:
		return m.handleTeamsKeys(msg)
	case ViewTeamMembers:
		return m.handleTeamMembersKeys(msg)
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
		m.searchTab = TabUsers
		m.searchInput.Focus()
		m.searchInput.SetValue("")
		m.searchInput.Placeholder = "Search GitHub users..."
		m.searchResults = nil
		m.lastQuery = ""
		m.err = nil
		m.updateSearchList() // Show recent/collaborators initially
		return m, nil

	case "t":
		m.view = ViewTeams
		m.loading = true
		m.err = nil
		m.searchInput.SetValue("")
		m.searchInput.Placeholder = "Filter teams..."
		m.searchInput.Blur()
		return m, loadTeams

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

	// Update text input and check for changes to trigger debounced search
	oldValue := m.searchInput.Value()
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	// Schedule debounced search if input changed
	newValue := m.searchInput.Value()
	if newValue != oldValue && m.searchInput.Focused() {
		m.debounceTimer++
		if newValue == "" {
			// Clear results immediately when input is empty
			m.searchResults = nil
			m.lastQuery = ""
			m.updateSearchList()
			return m, cmd
		}
		return m, tea.Batch(cmd, scheduleDebounce(newValue, m.debounceTimer))
	}

	return m, cmd
}

func (m Model) handleTeamsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if item, ok := m.teamList.SelectedItem().(teamItem); ok {
			m.selectedTeam = &item.team
			m.view = ViewTeamMembers
			m.loading = true
			m.searchInput.SetValue("")
			m.searchInput.Placeholder = "Filter team members..."
			m.searchInput.Blur() // Start with list focused, not input
			return m, loadTeamMembers(item.team.Org, item.team.Slug)
		}

	case "tab":
		if m.searchInput.Focused() {
			m.searchInput.Blur()
		} else {
			m.searchInput.Focus()
		}
		return m, nil

	case "up", "down":
		var cmd tea.Cmd
		m.teamList, cmd = m.teamList.Update(msg)
		return m, cmd
	}

	// Update text input and filter teams
	oldValue := m.searchInput.Value()
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	newValue := m.searchInput.Value()
	if newValue != oldValue {
		m.filterTeams(newValue)
		m.updateTeamList()
	}

	return m, cmd
}

func (m Model) handleTeamMembersKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// If input is focused and empty, blur it to allow list navigation
		if m.searchInput.Focused() && m.searchInput.Value() == "" {
			m.searchInput.Blur()
			return m, nil
		}
		// If input not focused, select from list
		if !m.searchInput.Focused() {
			if item, ok := m.searchList.SelectedItem().(pairItem); ok {
				if err := config.AddPair(item.pair); err != nil {
					m.err = err
					return m, nil
				}
				m.view = ViewMain
				m.selectedTeam = nil
				m.teamMembers = nil
				m.searchInput.SetValue("")
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
		// Always allow arrow key navigation in the list
		var cmd tea.Cmd
		m.searchList, cmd = m.searchList.Update(msg)
		return m, cmd
	}

	// Update text input and filter team members
	oldValue := m.searchInput.Value()
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	newValue := m.searchInput.Value()
	if newValue != oldValue {
		m.filterTeamMembers(newValue)
		m.updateSearchList()
	}

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

	// If viewing team members, show filtered team members
	if m.view == ViewTeamMembers {
		for _, p := range m.filteredTeamMembers {
			items = append(items, pairItem{pair: p})
		}
		m.searchList.SetItems(items)
		return
	}

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

func (m *Model) updateTeamList() {
	items := make([]list.Item, len(m.filteredTeams))
	for i, t := range m.filteredTeams {
		items[i] = teamItem{team: t}
	}
	m.teamList.SetItems(items)
}

func (m *Model) filterTeams(query string) {
	if query == "" {
		m.filteredTeams = m.teams
		return
	}

	query = strings.ToLower(query)
	filtered := make([]github.Team, 0)
	for _, t := range m.teams {
		if strings.Contains(strings.ToLower(t.Name), query) ||
			strings.Contains(strings.ToLower(t.Slug), query) ||
			strings.Contains(strings.ToLower(t.Org), query) {
			filtered = append(filtered, t)
		}
	}
	m.filteredTeams = filtered
}

func (m *Model) filterTeamMembers(query string) {
	if query == "" {
		m.filteredTeamMembers = m.teamMembers
		return
	}

	query = strings.ToLower(query)
	filtered := make([]config.Pair, 0)
	for _, p := range m.teamMembers {
		if strings.Contains(strings.ToLower(p.Username), query) ||
			strings.Contains(strings.ToLower(p.Name), query) {
			filtered = append(filtered, p)
		}
	}
	m.filteredTeamMembers = filtered
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
		return searchResultsMsg{results: results, query: query}
	}
}

func lookupUser(username string) tea.Cmd {
	return func() tea.Msg {
		pair, err := github.LookupUser(username)
		return userLookedUpMsg{pair: pair, err: err}
	}
}

func scheduleDebounce(query string, timerID int) tea.Cmd {
	return tea.Tick(debounceDelay, func(t time.Time) tea.Msg {
		return debounceTickMsg{query: query, timerID: timerID}
	})
}

func loadTeams() tea.Msg {
	teams, err := github.GetUserTeams()
	if err != nil {
		return errMsg{err: err}
	}
	return teamsLoadedMsg{teams: teams}
}

func loadTeamMembers(org, slug string) tea.Cmd {
	return func() tea.Msg {
		members, err := github.GetTeamMembers(org, slug)
		if err != nil {
			return errMsg{err: err}
		}
		return teamMembersLoadedMsg{members: members}
	}
}
