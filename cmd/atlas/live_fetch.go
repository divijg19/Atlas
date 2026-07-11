package main

import (
	"context"

	indexpkg "github.com/divijg19/Atlas/internal/index"
	"github.com/divijg19/Atlas/internal/live"
)

func buildLiveIndex(ctx context.Context, query string) (indexpkg.Index, error) {
	return live.BuildLiveIndex(ctx, query)
}

func fetchLiveUsernames(ctx context.Context, query string) ([]string, error) {
	return live.FetchLiveUsernames(ctx, query)
}
