# AGENTS.md - AI Assistant Guide for gh-pair

This document helps AI coding assistants understand and work effectively with the gh-pair codebase.

## Project Overview

**gh-pair** is a GitHub CLI extension that manages pair programming co-authors. It automatically adds `Co-Authored-By` trailers to git commit messages via a `prepare-commit-msg` hook.

### Tech Stack

- **Language**: Go 1.25.6
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm architecture)
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **GitHub API**: Via `gh` CLI (requires authenticated `gh`)

## Code Structure

```
gh-pair/
├── main.go                      # Entry point → cmd.Execute()
├── cmd/                         # Cobra CLI commands
│   ├── root.go                  # Root command (launches TUI when no args)
│   ├── add.go                   # gh pair add @username
│   ├── remove.go                # gh pair remove @username
│   ├── list.go                  # gh pair list
│   ├── clear.go                 # gh pair clear
│   └── init.go                  # gh pair init (installs git hook)
├── internal/
│   ├── config/config.go         # Pair CRUD, JSON persistence
│   ├── git/repo.go              # Git repository utilities
│   ├── github/client.go         # GitHub API via gh CLI
│   ├── hook/hook.go             # prepare-commit-msg hook management
│   └── tui/                     # Bubble Tea TUI
│       ├── app.go               # Run() entry point
│       ├── model.go             # Model struct, Init, Update
│       └── view.go              # View rendering
└── bin/                         # Build output directory
```

## Architecture Patterns

### Bubble Tea (Elm Architecture)

The TUI follows the Elm architecture pattern in `internal/tui/model.go`:

1. **Model** - Application state (`Model` struct)
2. **Update** - Message handling (`Update` method)
3. **View** - UI rendering (`View` method in `view.go`)

Key types:
- `View` enum: `ViewMain`, `ViewSearch`, `ViewTeams`, `ViewTeamMembers`, `ViewHelp`
- Messages: `pairsLoadedMsg`, `searchResultsMsg`, `userLookedUpMsg`, `errMsg`, etc.
- Commands return `tea.Cmd` for async operations

### GitHub API Pattern

All GitHub API calls go through `gh` CLI (not direct HTTP):

```go
cmd := exec.Command("gh", "api", "users/octocat")
output, err := cmd.Output()
```

This ensures authentication is handled by the user's `gh` session.

### Configuration Storage

Config is stored per-repository in `.git/gh-pair/`:
- `pairs.json` - Currently active pairs
- `recent.json` - Recently used pairs (max 10)

## Development Commands

```bash
# Build
go build -o gh-pair .

# Run the TUI
./gh-pair

# Run a subcommand
./gh-pair add @octocat
./gh-pair list

# Install as gh extension (for testing)
gh extension install .

# Run after local changes
gh extension remove pair && gh extension install .
```

## Key Conventions

### Error Handling

- Return errors up the call stack; don't log and return
- Use `fmt.Errorf("context: %w", err)` for wrapping
- TUI displays errors via `m.err` field

### Code Organization

- Public functions have doc comments
- Internal packages are in `internal/` (not importable externally)
- Each Cobra command is in its own file in `cmd/`

### TUI Patterns

- Key handlers are split by view: `handleMainKeys()`, `handleSearchKeys()`, etc.
- List components use `list.Model` from `bubbles`
- Debounced search: 300ms delay before API calls (`debounceDelay`)
- Spinner shown during loading states

### Naming

- `Pair` - A co-author (username, name, email)
- `pair` lowercase in variables, `Pair` for types
- Messages end in `Msg` suffix (e.g., `pairsLoadedMsg`)

## Important Files to Know

| File | Purpose |
|------|---------|
| `internal/tui/model.go` | Core TUI logic, state management |
| `internal/config/config.go` | Pair storage/retrieval |
| `internal/hook/hook.go` | Git hook script (embedded as string) |
| `internal/github/client.go` | All GitHub API interactions |
| `cmd/root.go` | CLI entry, launches TUI |

## Testing Locally

1. Ensure `gh` is authenticated: `gh auth status`
2. Run from a git repository (required)
3. Build and run: `go build -o gh-pair . && ./gh-pair`

## Gotchas

- Must be run inside a git repository (checks `git rev-parse --is-inside-work-tree`)
- GitHub API requires `gh` CLI to be installed and authenticated
- The `prepare-commit-msg` hook uses shell script (not Go) for portability
- Team features require org membership and appropriate permissions
