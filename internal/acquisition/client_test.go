package acquisition

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func TestFetchReposSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/octocat/repos" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"name":"repo1","fork":false,"size":120,"updated_at":"2024-01-02T03:04:05Z"}]`))
	}))
	defer srv.Close()

	repos, err := testClient(srv.URL).FetchRepos(context.Background(), "octocat")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(repos))
	}
	if repos[0].Name != "repo1" || repos[0].Fork || repos[0].Size != 120 {
		t.Errorf("unexpected repo dto: %+v", repos[0])
	}
	if repos[0].UpdatedAt != "2024-01-02T03:04:05Z" {
		t.Errorf("expected raw timestamp string, got %q", repos[0].UpdatedAt)
	}
}

func TestFetchReposEmptyUsername(t *testing.T) {
	_, err := testClient("http://example.invalid").FetchRepos(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty username")
	}
}

func TestFetchReposAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer srv.Close()

	_, err := testClient(srv.URL).FetchRepos(context.Background(), "ghost")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "Not Found" {
		t.Errorf("expected message %q, got %q", "Not Found", apiErr.Message)
	}
}

func TestFetchUserSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/octocat" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"The Octocat","company":"GitHub","location":"SF","followers":10,"following":2,"created_at":"2011-01-25T18:44:36Z"}`))
	}))
	defer srv.Close()

	user, err := testClient(srv.URL).FetchUser(context.Background(), "octocat")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "The Octocat" || user.Company != "GitHub" || user.Location != "SF" {
		t.Errorf("unexpected user dto: %+v", user)
	}
	if user.CreatedAt != "2011-01-25T18:44:36Z" {
		t.Errorf("expected raw created_at string, got %q", user.CreatedAt)
	}
}

func TestFetchUserAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"message":"rate limit"}`))
	}))
	defer srv.Close()

	_, err := testClient(srv.URL).FetchUser(context.Background(), "octocat")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "status 403") || !strings.Contains(err.Error(), "rate limit") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestFetchContributionsSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		w.WriteHeader(http.StatusOK)
		if strings.Contains(q, "type:pr") {
			_, _ = w.Write([]byte(`{"total_count":5}`))
			return
		}
		_, _ = w.Write([]byte(`{"total_count":3}`))
	}))
	defer srv.Close()

	dto, err := testClient(srv.URL).FetchContributions(context.Background(), "octocat")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.PullRequests != 5 {
		t.Errorf("expected 5 PRs, got %d", dto.PullRequests)
	}
	if dto.Issues != 3 {
		t.Errorf("expected 3 issues, got %d", dto.Issues)
	}
}

func TestFetchContributionsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"message":"bad query"}`))
	}))
	defer srv.Close()

	_, err := testClient(srv.URL).FetchContributions(context.Background(), "octocat")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "fetch pull requests") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestSearchRepositoryOwnersDedup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/repositories" {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"items":[{"owner":{"login":"alice"}},{"owner":{"login":"bob"}},{"owner":{"login":"alice"}},{"owner":{"login":""}}]}`))
	}))
	defer srv.Close()

	owners, err := testClient(srv.URL).SearchRepositoryOwners(context.Background(), "backend")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(owners) != 2 || owners[0] != "alice" || owners[1] != "bob" {
		t.Errorf("unexpected owners: %v", owners)
	}
}

func TestSearchRepositoryOwnersEmptyQuery(t *testing.T) {
	owners, err := testClient("http://example.invalid").SearchRepositoryOwners(context.Background(), "   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owners != nil {
		t.Errorf("expected nil owners for empty query, got %v", owners)
	}
}
