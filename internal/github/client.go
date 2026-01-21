package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/omgitsads/gh-pair/internal/config"
)

// userResponse represents the GitHub API response for a user.
type userResponse struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    int    `json:"id"`
}

// searchResponse represents the GitHub API response for user search.
type searchResponse struct {
	Items []userResponse `json:"items"`
}

// LookupUser fetches a GitHub user by username and returns a Pair.
func LookupUser(username string) (*config.Pair, error) {
	// Strip @ prefix if present
	username = strings.TrimPrefix(username, "@")

	cmd := exec.Command("gh", "api", fmt.Sprintf("users/%s", username))
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("user not found: %s (gh api error: %s)", username, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to lookup user: %w", err)
	}

	var user userResponse
	if err := json.Unmarshal(output, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	// Use noreply email if user has no public email
	email := user.Email
	if email == "" {
		email = fmt.Sprintf("%d+%s@users.noreply.github.com", user.ID, user.Login)
	}

	// Use login as name if name is empty
	name := user.Name
	if name == "" {
		name = user.Login
	}

	return &config.Pair{
		Username: user.Login,
		Name:     name,
		Email:    email,
	}, nil
}

// SearchUsers searches for GitHub users matching the query.
func SearchUsers(query string) ([]config.Pair, error) {
	if query == "" {
		return []config.Pair{}, nil
	}

	cmd := exec.Command("gh", "api", fmt.Sprintf("search/users?q=%s&per_page=10", query))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	var response searchResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	pairs := make([]config.Pair, 0, len(response.Items))
	for _, user := range response.Items {
		// For search results, we only have basic info
		// Full details would require another API call per user
		email := fmt.Sprintf("%d+%s@users.noreply.github.com", user.ID, user.Login)
		name := user.Login
		if user.Name != "" {
			name = user.Name
		}

		pairs = append(pairs, config.Pair{
			Username: user.Login,
			Name:     name,
			Email:    email,
		})
	}

	return pairs, nil
}

// GetRepoCollaborators fetches collaborators for the current repository.
func GetRepoCollaborators() ([]config.Pair, error) {
	// Get the current repo info
	cmd := exec.Command("gh", "repo", "view", "--json", "owner,name")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get repo info: %w", err)
	}

	var repoInfo struct {
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(output, &repoInfo); err != nil {
		return nil, fmt.Errorf("failed to parse repo info: %w", err)
	}

	// Fetch collaborators
	cmd = exec.Command("gh", "api", fmt.Sprintf("repos/%s/%s/collaborators?per_page=20", repoInfo.Owner.Login, repoInfo.Name))
	output, err = cmd.Output()
	if err != nil {
		// User might not have permission to list collaborators
		return []config.Pair{}, nil
	}

	var collaborators []userResponse
	if err := json.Unmarshal(output, &collaborators); err != nil {
		return nil, fmt.Errorf("failed to parse collaborators: %w", err)
	}

	pairs := make([]config.Pair, 0, len(collaborators))
	for _, user := range collaborators {
		email := fmt.Sprintf("%d+%s@users.noreply.github.com", user.ID, user.Login)
		name := user.Login
		if user.Name != "" {
			name = user.Name
		}

		pairs = append(pairs, config.Pair{
			Username: user.Login,
			Name:     name,
			Email:    email,
		})
	}

	return pairs, nil
}
