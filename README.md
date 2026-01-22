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

1. Run `gh pair init` in your repository to install the `commit-msg` hook
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

## Themes

gh-pair supports multiple color themes. Use `--theme` to override temporarily:

```bash
gh pair --theme dracula
gh pair --theme nord
```

### Setting a Default Theme

```bash
gh pair theme set dracula
```

This saves your preference to the global config file.

### Available Themes

| Theme | Description |
|-------|-------------|
| `default` | Blue and magenta (default) |
| `dracula` | Purple/pink dark theme |
| `nord` | Blue arctic palette |
| `solarized-dark` | Warm dark theme |
| `solarized-light` | Warm light theme |
| `catppuccin` | Pastel dark theme |

List all themes: `gh pair theme`
Preview a theme: `gh pair theme dracula`

### Custom Themes

Create custom themes by adding JSON files to `~/.config/gh-pair/themes/`:

Example theme file:

```json
{
  "name": "my-theme",
  "colors": {
    "primary": "#bd93f9",
    "secondary": "#ff79c6",
    "success": "#50fa7b",
    "error": "#ff5555",
    "warning": "#ffb86c",
    "text": "#f8f8f2",
    "textDim": "#6272a4",
    "border": "#bd93f9",
    "accent": "#ff79c6"
  }
}
```

Save as `my-theme.json` in the themes directory, then use with `gh pair theme set my-theme`.

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
