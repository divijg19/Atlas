package main

import (
	"context"

	"github.com/divijg19/GH-Analyzer/internal/live"
	indexpkg "github.com/divijg19/GH-Analyzer/internal/index"
)

func buildLiveIndex(ctx context.Context, query string) (indexpkg.Index, error) {
	return live.BuildLiveIndex(ctx, query)
}

func fetchLiveUsernames(ctx context.Context, query string) ([]string, error) {
	return live.FetchLiveUsernames(ctx, query)
}
