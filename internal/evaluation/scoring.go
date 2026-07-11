package evaluation

import (
	"github.com/divijg19/Atlas/internal/index"
	"github.com/divijg19/Atlas/internal/signals"
)

const (
	ownershipWeight   = 0.3
	consistencyWeight = 0.4
	depthWeight       = 0.3

	smallSampleRepoThreshold     = 3
	smallSamplePenaltyMultiplier = 7
)

// defaultRankingWeights returns the canonical indicator weights used for both
// search ranking and overall scoring. Evaluation is the single owner of these
// weights.
func defaultRankingWeights() map[string]float64 {
	return map[string]float64{
		signals.SignalConsistency: consistencyWeight,
		signals.SignalOwnership:   ownershipWeight,
		signals.SignalDepth:       depthWeight,
	}
}

// RankingPolicy computes a weighted score for a candidate profile from its
// indicator signals. It is the concrete ranking implementation owned by the
// evaluation layer.
type RankingPolicy struct {
	Weights map[string]float64
}

// Score returns the weighted sum of a profile's indicator signals.
func (w RankingPolicy) Score(p index.Profile) float64 {
	weights := w.Weights
	if len(weights) == 0 {
		weights = defaultRankingWeights()
	}

	score := 0.0
	for signal, weight := range weights {
		value, ok := p.Signals[signal]
		if !ok {
			continue
		}

		score += value * weight
	}

	return score
}

// OverallScore computes the weighted overall score (0-100) from raw component
// scores.
func OverallScore(rs signals.RawScore) int {
	return int(float64(rs.Consistency)*consistencyWeight +
		float64(rs.Ownership)*ownershipWeight +
		float64(rs.Depth)*depthWeight)
}

// ApplySmallSamplePenalty applies the small sample penalty to a raw overall
// score. Profiles with fewer than the threshold repository count are scaled
// down to reflect reduced confidence.
func ApplySmallSamplePenalty(rawScore int, repoCount int) int {
	if repoCount >= smallSampleRepoThreshold {
		return rawScore
	}

	return (rawScore * smallSamplePenaltyMultiplier) / 10
}
