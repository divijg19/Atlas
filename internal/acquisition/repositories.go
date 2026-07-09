package acquisition

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/divijg19/GH-Analyzer/internal/signals"
)

// RepoDTO mirrors a repository object from the GitHub REST API. Field names and
// types follow GitHub's schema; timestamps remain strings and are parsed during
// normalization.
type RepoDTO struct {
	Name      string `json:"name"`
	Fork      bool   `json:"fork"`
	Size      int    `json:"size"`
	UpdatedAt string `json:"updated_at"`
}

// FetchRepos retrieves a user's public repositories from the GitHub REST API.
func (c *Client) FetchRepos(ctx context.Context, username string) ([]RepoDTO, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	resp, err := c.get(ctx, fmt.Sprintf("%s/users/%s/repos", c.baseURL, username))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var githubError struct {
			Message string `json:"message"`
		}

		message := resp.Status
		if err := json.NewDecoder(resp.Body).Decode(&githubError); err == nil && strings.TrimSpace(githubError.Message) != "" {
			message = strings.TrimSpace(githubError.Message)
		}

		return nil, APIError{
			StatusCode: resp.StatusCode,
			Message:    message,
		}
	}

	var repos []RepoDTO
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

// FetchReposNormalized retrieves a user's public repositories and returns them
// as domain models. It is equivalent to FetchRepos followed by NormalizeRepos.
func (c *Client) FetchReposNormalized(ctx context.Context, username string) ([]signals.Repo, error) {
	dtos, err := c.FetchRepos(ctx, username)
	if err != nil {
		return nil, err
	}

	return NormalizeRepos(dtos), nil
}
