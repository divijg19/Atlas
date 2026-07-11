package acquisition

import (
	"testing"
)

func TestNormalizeRepo(t *testing.T) {
	repo := normalizeRepo(RepoDTO{
		Name:      "repo1",
		Fork:      false,
		Size:      120,
		UpdatedAt: "2024-01-02T03:04:05Z",
	})

	if repo.Name != "repo1" || repo.Fork || repo.Size != 120 {
		t.Errorf("unexpected repo: %+v", repo)
	}
	if repo.UpdatedAt.IsZero() {
		t.Errorf("expected parsed updated_at, got zero time")
	}
}

func TestNormalizeRepoEmptyTime(t *testing.T) {
	repo := normalizeRepo(RepoDTO{Name: "x", UpdatedAt: ""})
	if !repo.UpdatedAt.IsZero() {
		t.Errorf("expected zero time for empty timestamp, got %v", repo.UpdatedAt)
	}
}

func TestNormalizeRepos(t *testing.T) {
	repos := NormalizeRepos([]RepoDTO{
		{Name: "a", Size: 1, UpdatedAt: "2024-01-01T00:00:00Z"},
		{Name: "b", Fork: true, Size: 2, UpdatedAt: "2024-02-01T00:00:00Z"},
	})
	if len(repos) != 2 {
		t.Fatalf("expected 2 repos, got %d", len(repos))
	}
	if repos[0].Name != "a" || repos[1].Fork != true {
		t.Errorf("unexpected repos: %+v", repos)
	}
}

func TestNormalizeUser(t *testing.T) {
	user := NormalizeUser(&UserDTO{
		Name:      "The Octocat",
		Bio:       "cat",
		Location:  "SF",
		Company:   "GitHub",
		Followers: 10,
		Following: 2,
		CreatedAt: "2011-01-25T18:44:36Z",
	})

	if user == nil {
		t.Fatal("expected non-nil metadata")
	}
	if user.Name != "The Octocat" || user.Company != "GitHub" || user.Location != "SF" {
		t.Errorf("unexpected metadata: %+v", user)
	}
	if user.Followers != 10 || user.Following != 2 {
		t.Errorf("unexpected follower counts: %+v", user)
	}
	if user.CreatedAt.IsZero() {
		t.Errorf("expected parsed created_at, got zero time")
	}
}

func TestNormalizeUserNil(t *testing.T) {
	if NormalizeUser(nil) != nil {
		t.Error("expected nil for nil DTO")
	}
}

func TestNormalizeUserZeroTime(t *testing.T) {
	user := NormalizeUser(&UserDTO{Name: "x", CreatedAt: ""})
	if !user.CreatedAt.IsZero() {
		t.Errorf("expected zero time for empty created_at, got %v", user.CreatedAt)
	}
}

func TestNormalizeContributions(t *testing.T) {
	summary := NormalizeContributions(&ContributionsDTO{PullRequests: 5, Issues: 3})

	if summary == nil {
		t.Fatal("expected non-nil summary")
	}
	if summary.TotalPullRequests != 5 {
		t.Errorf("expected 5 PRs, got %d", summary.TotalPullRequests)
	}
	if summary.IssuesOpened != 3 {
		t.Errorf("expected 3 issues, got %d", summary.IssuesOpened)
	}
	if summary.TotalContributions != 8 {
		t.Errorf("expected total 8, got %d", summary.TotalContributions)
	}
}

func TestNormalizeContributionsNil(t *testing.T) {
	if NormalizeContributions(nil) != nil {
		t.Error("expected nil for nil DTO")
	}
}
