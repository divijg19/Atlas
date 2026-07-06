package live

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/divijg19/GH-Analyzer/internal/contributions"
	"github.com/divijg19/GH-Analyzer/internal/github"
	"github.com/divijg19/GH-Analyzer/internal/signals"
)

const (
	MaxRepoResults = 30
	MaxUsers       = 20
)

var RepoSearchURL = "https://api.github.com/search/repositories"

type repoSearchResponse struct {
	Items []repoItem `json:"items"`
}

type repoItem struct {
	Owner repoOwner `json:"owner"`
}

type repoOwner struct {
	Login string `json:"login"`
}

// FetchLiveUsernames searches GitHub repositories and returns deduplicated
// owner usernames, up to MaxUsers.
func FetchLiveUsernames(ctx context.Context, query string) ([]string, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return nil, nil
	}

	endpoint := fmt.Sprintf("%s?q=%s&per_page=%d", RepoSearchURL, url.QueryEscape(trimmed), MaxRepoResults)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub data")
	}

	resp, err := github.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub data")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch GitHub data")
	}

	var payload repoSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub data")
	}

	seen := make(map[string]struct{}, len(payload.Items))
	usernames := make([]string, 0, MaxUsers)

	for _, item := range payload.Items {
		login := strings.TrimSpace(item.Owner.Login)
		if login == "" {
			continue
		}
		if _, ok := seen[login]; ok {
			continue
		}
		seen[login] = struct{}{}
		usernames = append(usernames, login)
		if len(usernames) >= MaxUsers {
			break
		}
	}

	return usernames, nil
}

// These function variables are overridable for testing.
var (
	FetchReposFn         = signals.FetchRepos
	FetchContributionsFn = contributions.FetchContributions
)
