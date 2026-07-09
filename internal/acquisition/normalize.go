package acquisition

import (
	"time"

	"github.com/divijg19/GH-Analyzer/internal/contributions"
	"github.com/divijg19/GH-Analyzer/internal/profile"
	"github.com/divijg19/GH-Analyzer/internal/signals"
)

// RepoDTO maps a GitHub repository DTO to the domain repository model.
func NormalizeRepo(dto RepoDTO) signals.Repo {
	return signals.Repo{
		Name:      dto.Name,
		Fork:      dto.Fork,
		Size:      dto.Size,
		UpdatedAt: parseTime(dto.UpdatedAt),
	}
}

// NormalizeRepos maps a slice of repository DTOs to domain repositories.
func NormalizeRepos(dtos []RepoDTO) []signals.Repo {
	repos := make([]signals.Repo, len(dtos))
	for i, dto := range dtos {
		repos[i] = NormalizeRepo(dto)
	}
	return repos
}

// NormalizeUser maps a GitHub user DTO to the domain metadata model.
func NormalizeUser(dto *UserDTO) *profile.UserMetadata {
	if dto == nil {
		return nil
	}

	return &profile.UserMetadata{
		Name:      dto.Name,
		Bio:       dto.Bio,
		Location:  dto.Location,
		Company:   dto.Company,
		Followers: dto.Followers,
		Following: dto.Following,
		CreatedAt: parseTime(dto.CreatedAt),
	}
}

// NormalizeContributions maps the raw contribution DTO to the domain summary.
// Total contributions aggregate both pull requests and issues.
func NormalizeContributions(dto *ContributionsDTO) *contributions.Summary {
	if dto == nil {
		return nil
	}

	return &contributions.Summary{
		TotalContributions: dto.PullRequests + dto.Issues,
		TotalPullRequests:  dto.PullRequests,
		IssuesOpened:       dto.Issues,
	}
}

// parseTime parses an RFC3339 timestamp. An empty or unparseable value yields
// the zero time, mirroring prior decoding behavior for optional timestamps.
func parseTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}

	return parsed
}
