package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/omgitsads/gh-pair/internal/git"
)

const (
	PairsFileName  = "pairs.json"
	RecentFileName = "recent.json"
	MaxRecentPairs = 10
)

// Pair represents a co-author for commits.
type Pair struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// CoAuthorLine returns the Co-Authored-By trailer line for this pair.
func (p Pair) CoAuthorLine() string {
	return "Co-Authored-By: " + p.Name + " <" + p.Email + ">"
}

// PairsConfig holds the current active pairs.
type PairsConfig struct {
	Pairs []Pair `json:"pairs"`
}

// RecentConfig holds recently used pairs for quick access.
type RecentConfig struct {
	Recent []Pair `json:"recent"`
}

// LoadPairs loads the current pairs from the config file.
func LoadPairs() (*PairsConfig, error) {
	configDir, err := git.ConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, PairsFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &PairsConfig{Pairs: []Pair{}}, nil
		}
		return nil, err
	}

	var config PairsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SavePairs saves the pairs configuration to the config file.
func SavePairs(config *PairsConfig) error {
	configDir, err := git.EnsureConfigDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(configDir, PairsFileName)
	return os.WriteFile(path, data, 0644)
}

// AddPair adds a pair to the config if not already present.
func AddPair(pair Pair) error {
	config, err := LoadPairs()
	if err != nil {
		return err
	}

	// Check if already exists
	for _, p := range config.Pairs {
		if p.Username == pair.Username {
			return nil // Already exists
		}
	}

	config.Pairs = append(config.Pairs, pair)

	if err := SavePairs(config); err != nil {
		return err
	}

	// Also add to recent
	return AddToRecent(pair)
}

// RemovePair removes a pair from the config by username.
func RemovePair(username string) error {
	config, err := LoadPairs()
	if err != nil {
		return err
	}

	newPairs := make([]Pair, 0, len(config.Pairs))
	for _, p := range config.Pairs {
		if p.Username != username {
			newPairs = append(newPairs, p)
		}
	}

	config.Pairs = newPairs
	return SavePairs(config)
}

// ClearPairs removes all pairs from the config.
func ClearPairs() error {
	config := &PairsConfig{Pairs: []Pair{}}
	return SavePairs(config)
}

// LoadRecent loads the recent pairs from the config file.
func LoadRecent() (*RecentConfig, error) {
	configDir, err := git.ConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, RecentFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &RecentConfig{Recent: []Pair{}}, nil
		}
		return nil, err
	}

	var config RecentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveRecent saves the recent pairs configuration.
func SaveRecent(config *RecentConfig) error {
	configDir, err := git.EnsureConfigDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(configDir, RecentFileName)
	return os.WriteFile(path, data, 0644)
}

// AddToRecent adds a pair to the recent list (moves to front if exists).
func AddToRecent(pair Pair) error {
	config, err := LoadRecent()
	if err != nil {
		return err
	}

	// Remove if already exists
	newRecent := make([]Pair, 0, MaxRecentPairs)
	for _, p := range config.Recent {
		if p.Username != pair.Username {
			newRecent = append(newRecent, p)
		}
	}

	// Add to front
	newRecent = append([]Pair{pair}, newRecent...)

	// Trim to max size
	if len(newRecent) > MaxRecentPairs {
		newRecent = newRecent[:MaxRecentPairs]
	}

	config.Recent = newRecent
	return SaveRecent(config)
}
