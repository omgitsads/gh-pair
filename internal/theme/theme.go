// Package theme provides theming support for the gh-pair TUI.
package theme

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Theme defines the color palette for the TUI.
type Theme struct {
	Name      string `json:"name"`
	Colors    Colors `json:"colors"`
}

// Colors defines the color values used throughout the TUI.
type Colors struct {
	Primary   string `json:"primary"`   // Main accent color (titles, borders)
	Secondary string `json:"secondary"` // Secondary accent (help keys)
	Success   string `json:"success"`   // Success messages
	Error     string `json:"error"`     // Error messages
	Warning   string `json:"warning"`   // Warning messages
	Text      string `json:"text"`      // Primary text
	TextDim   string `json:"textDim"`   // Dimmed/subtle text
	Border    string `json:"border"`    // Border color
	Accent    string `json:"accent"`    // Accent highlights
}

// Preset themes
var (
	DefaultTheme = Theme{
		Name: "default",
		Colors: Colors{
			Primary:   "99",  // Blue
			Secondary: "205", // Magenta
			Success:   "42",  // Green
			Error:     "196", // Red
			Warning:   "214", // Orange
			Text:      "252", // Light gray
			TextDim:   "241", // Dark gray
			Border:    "99",  // Blue
			Accent:    "205", // Magenta
		},
	}

	DraculaTheme = Theme{
		Name: "dracula",
		Colors: Colors{
			Primary:   "#bd93f9", // Purple
			Secondary: "#ff79c6", // Pink
			Success:   "#50fa7b", // Green
			Error:     "#ff5555", // Red
			Warning:   "#ffb86c", // Orange
			Text:      "#f8f8f2", // Foreground
			TextDim:   "#6272a4", // Comment
			Border:    "#bd93f9", // Purple
			Accent:    "#ff79c6", // Pink
		},
	}

	NordTheme = Theme{
		Name: "nord",
		Colors: Colors{
			Primary:   "#88c0d0", // Frost cyan
			Secondary: "#81a1c1", // Frost blue
			Success:   "#a3be8c", // Aurora green
			Error:     "#bf616a", // Aurora red
			Warning:   "#ebcb8b", // Aurora yellow
			Text:      "#eceff4", // Snow Storm
			TextDim:   "#4c566a", // Polar Night
			Border:    "#5e81ac", // Frost dark blue
			Accent:    "#b48ead", // Aurora purple
		},
	}

	SolarizedDarkTheme = Theme{
		Name: "solarized-dark",
		Colors: Colors{
			Primary:   "#268bd2", // blue
			Secondary: "#2aa198", // cyan
			Success:   "#859900", // green
			Error:     "#dc322f", // red
			Warning:   "#b58900", // yellow
			Text:      "#839496", // base0 (light text on dark bg)
			TextDim:   "#586e75", // base01 (dimmed on dark bg)
			Border:    "#6c71c4", // violet
			Accent:    "#d33682", // magenta
		},
	}

	SolarizedLightTheme = Theme{
		Name: "solarized-light",
		Colors: Colors{
			Primary:   "#268bd2", // blue
			Secondary: "#2aa198", // cyan
			Success:   "#859900", // green
			Error:     "#dc322f", // red
			Warning:   "#b58900", // yellow
			Text:      "#657b83", // base00 (dark text on light bg)
			TextDim:   "#93a1a1", // base1 (dimmed on light bg)
			Border:    "#6c71c4", // violet
			Accent:    "#d33682", // magenta
		},
	}

	CatppuccinTheme = Theme{
		Name: "catppuccin",
		Colors: Colors{
			Primary:   "#cba6f7", // Mauve
			Secondary: "#f5c2e7", // Pink
			Success:   "#a6e3a1", // Green
			Error:     "#f38ba8", // Red
			Warning:   "#f9e2af", // Yellow
			Text:      "#cdd6f4", // Text
			TextDim:   "#6c7086", // Overlay0
			Border:    "#89b4fa", // Blue
			Accent:    "#f5c2e7", // Pink
		},
	}
)

// presetThemes maps theme names to their definitions.
var presetThemes = map[string]Theme{
	"default":         DefaultTheme,
	"dracula":         DraculaTheme,
	"nord":            NordTheme,
	"solarized-dark":  SolarizedDarkTheme,
	"solarized-light": SolarizedLightTheme,
	"catppuccin":      CatppuccinTheme,
}

// PresetNames returns a list of all preset theme names.
func PresetNames() []string {
	return []string{"default", "dracula", "nord", "solarized-dark", "solarized-light", "catppuccin"}
}

// GetTheme returns a theme by name. It checks custom themes first,
// then falls back to presets. Returns default theme if not found.
func GetTheme(name string) Theme {
	if name == "" {
		name = "default"
	}

	// Try loading custom theme from config dir
	if theme, err := loadCustomTheme(name); err == nil {
		return theme
	}

	// Fall back to preset
	if theme, ok := presetThemes[name]; ok {
		return theme
	}

	return DefaultTheme
}

// loadCustomTheme loads a theme from ~/.config/gh-pair/themes/{name}.json
func loadCustomTheme(name string) (Theme, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Theme{}, err
	}

	themePath := filepath.Join(home, ".config", "gh-pair", "themes", name+".json")
	data, err := os.ReadFile(themePath)
	if err != nil {
		return Theme{}, err
	}

	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return Theme{}, err
	}

	return theme, nil
}

// ListCustomThemes returns names of custom themes in ~/.config/gh-pair/themes/.
func ListCustomThemes() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	themesDir := filepath.Join(home, ".config", "gh-pair", "themes")
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return nil
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			name := entry.Name()[:len(entry.Name())-5] // Remove .json
			names = append(names, name)
		}
	}
	return names
}
