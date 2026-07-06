package live

import (
	"context"
	"fmt"

	indexpkg "github.com/divijg19/GH-Analyzer/internal/index"
	"github.com/divijg19/GH-Analyzer/internal/profile"
	"github.com/divijg19/GH-Analyzer/internal/signals"
)

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
		profile, err := buildLiveProfile(ctx, username)
		if err != nil {
			continue
		}
		idx.Add(profile)
	}

	return idx, nil
}

func buildLiveProfile(ctx context.Context, username string) (indexpkg.Profile, error) {
	repos, err := FetchReposFn(ctx, username)
	if err != nil {
		return indexpkg.Profile{}, fmt.Errorf("fetch repos for %q: %w", username, err)
	}

	facts := signals.FromRepos(repos)

	meta, err := profile.FetchUserMetadata(ctx, username)
	if err != nil {
		return indexpkg.Profile{}, fmt.Errorf("fetch metadata for %q: %w", username, err)
	}

	contribSummary, err := FetchContributionsFn(ctx, username)
	if err != nil {
		return indexpkg.Profile{}, fmt.Errorf("fetch contributions for %q: %w", username, err)
	}

	signalValues := signals.ExtractSignalsFromFacts(facts)
	scores := signals.ScoreSignals(signalValues)
	report := signals.BuildReport(username, scores, repos)

	return indexpkg.Profile{
		Username:      username,
		Signals:       signals.SignalsFromReport(report),
		Facts:         &facts,
		Metadata:      meta,
		Contributions: contribSummary,
	}, nil
}
