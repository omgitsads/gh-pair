package theme

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the global gh-pair configuration.
type Config struct {
	Theme string `json:"theme"`
}

// configDir returns the path to the gh-pair config directory (~/.config/gh-pair).
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gh-pair"), nil
}

// configPath returns the path to the global config file.
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// LoadConfig loads the global configuration from ~/.config/gh-pair/config.json.
// Returns default config if file doesn't exist.
func LoadConfig() Config {
	path, err := configPath()
	if err != nil {
		return Config{Theme: "default"}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{Theme: "default"}
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{Theme: "default"}
	}

	if cfg.Theme == "" {
		cfg.Theme = "default"
	}

	return cfg
}

// SaveConfig saves the global configuration to ~/.config/gh-pair/config.json.
func SaveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetConfiguredTheme returns the theme name from global config.
func GetConfiguredTheme() string {
	return LoadConfig().Theme
}

// SetConfiguredTheme saves the theme name to global config.
func SetConfiguredTheme(themeName string) error {
	cfg := LoadConfig()
	cfg.Theme = themeName
	return SaveConfig(cfg)
}
