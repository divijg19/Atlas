package engine

import (
	"math"
	"testing"

	"github.com/divijg19/Atlas/internal/evaluation"
	"github.com/divijg19/Atlas/internal/index"
)

func TestWeightedRankingSkipsMissingSignals(t *testing.T) {
	ranking := evaluation.RankingPolicy{}
	profile := index.Profile{Signals: map[string]float64{
		"consistency": 0.8,
		// ownership is intentionally missing
		"depth": 0.5,
	}}

	got := ranking.Score(profile)
	want := (0.8 * 0.4) + (0.5 * 0.3)
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("expected %.2f, got %.2f", want, got)
	}
}
