package github

import (
	"context"
	"net/http"
	"os"
	"time"
)

const (
	defaultUserAgent = "gh-analyzer"
	defaultTimeout   = 30 * time.Second
)

var sharedClient = &http.Client{Timeout: defaultTimeout}

func Client() *http.Client {
	return sharedClient
}

func AuthHeader() (string, string) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", ""
	}
	return "Authorization", "Bearer " + token
}

func SetHeaders(req *http.Request) {
	req.Header.Set("User-Agent", defaultUserAgent)
	if key, val := AuthHeader(); key != "" {
		req.Header.Set(key, val)
	}
}

func Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	SetHeaders(req)
	return Client().Do(req)
}
