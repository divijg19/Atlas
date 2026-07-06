package main

import (
	"context"

	"github.com/divijg19/GH-Analyzer/internal/live"
	indexpkg "github.com/divijg19/GH-Analyzer/internal/index"
)

func buildLiveIndexForServer(ctx context.Context, query string) (indexpkg.Index, error) {
	return live.BuildLiveIndex(ctx, query)
}
