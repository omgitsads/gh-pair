# gh-pair

A GitHub CLI extension for managing pair programming co-authors. Automatically adds `Co-Authored-By` trailers to your commit messages.

## Installation

```bash
gh extension install omgitsads/gh-pair
```

## Usage

### Interactive TUI

Launch the interactive interface to manage your pairs:

```bash
gh pair
```

### Quick Commands

```bash
# Install the git hook in your repository
gh pair init

# Add a pair by GitHub username
gh pair add @octocat

# Remove a pair
gh pair remove @octocat

# List current pairs
gh pair list

# Clear all pairs
gh pair clear
```

## How It Works

1. Run `gh pair init` in your repository to install the `prepare-commit-msg` hook
2. Add collaborators you're pairing with using `gh pair add @username` or the TUI
3. Your commits will automatically include `Co-Authored-By` trailers

### Example

```bash
$ gh pair add @octocat
Added: The Octocat <octocat@github.com>

$ git commit -m "Add new feature"
# Commit message becomes:
# Co-Authored-By: The Octocat <octocat@github.com>
#
# Add new feature
```

## Configuration

Configuration is stored per-repository in `.git/gh-pair/`:
- `pairs.json` - Current active pairs
- `recent.json` - Recently used pairs for quick access

## Requirements

- [GitHub CLI](https://cli.github.com/) (`gh`) installed and authenticated
- Git repository

## Keyboard Shortcuts (TUI)

| Key | Action |
|-----|--------|
| `a` | Add a new pair |
| `d` / `Delete` | Remove selected pair |
| `c` | Clear all pairs |
| `/` | Search GitHub users |
| `↑` / `↓` | Navigate list |
| `Enter` | Select / Confirm |
| `Esc` | Cancel / Back |
| `?` | Show help |
| `q` | Quit |

## License

MIT
