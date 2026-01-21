package git

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrNotARepository = errors.New("not a git repository")
)

// RepoRoot returns the root directory of the current git repository.
func RepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", ErrNotARepository
	}
	return strings.TrimSpace(string(output)), nil
}

// GitDir returns the .git directory path for the current repository.
func GitDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", ErrNotARepository
	}
	gitDir := strings.TrimSpace(string(output))

	// Convert relative path to absolute
	if !filepath.IsAbs(gitDir) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		gitDir = filepath.Join(cwd, gitDir)
	}

	return gitDir, nil
}

// HooksDir returns the path to the hooks directory.
func HooksDir() (string, error) {
	gitDir, err := GitDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(gitDir, "hooks"), nil
}

// ConfigDir returns the path to the gh-pair config directory within .git.
func ConfigDir() (string, error) {
	gitDir, err := GitDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(gitDir, "gh-pair"), nil
}

// EnsureConfigDir creates the gh-pair config directory if it doesn't exist.
func EnsureConfigDir() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return configDir, nil
}

// IsInsideWorkTree checks if we're inside a git work tree.
func IsInsideWorkTree() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}
