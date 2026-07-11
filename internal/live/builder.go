// Package live builds candidate profiles by fetching directly from GitHub.
//
// It composes the acquisition, signals, and index layers to produce a Profile
// on demand (the "live" path), as opposed to reading from a stored dataset.
// See docs/INTELLIGENCE.md.
package live

import (
	"context"

	"github.com/divijg19/Atlas/internal/acquisition"
	indexpkg "github.com/divijg19/Atlas/internal/index"
)

// fetcher is the acquisition surface the live pipeline depends on. It is
// satisfied by *acquisition.Client and overridable in tests.
type fetcher interface {
	indexpkg.ProfileFetcher
	SearchRepositoryOwners(ctx context.Context, query string) ([]string, error)
}

// Client is the acquisition client used for live fetching. Tests may override
// it with a fetcher implementation; production code uses the default
// GitHub-backed client.
var Client fetcher = acquisition.NewClient()

// BuildLiveIndex fetches live data from GitHub for a search query and
// returns a complete index with profiles containing signals, facts,
// metadata, and contributions.
func BuildLiveIndex(ctx context.Context, query string) (indexpkg.Index, error) {
	usernames, err := FetchLiveUsernames(ctx, query)
	if err != nil {
		return indexpkg.Index{}, err
	}

	idx := indexpkg.Index{Profiles: make([]indexpkg.Profile, 0, len(usernames))}
	for _, username := range usernames {
		profile, err := indexpkg.BuildProfile(ctx, Client, username)
		if err != nil {
			continue
		}
		idx.Add(profile)
	}

	return idx, nil
}
