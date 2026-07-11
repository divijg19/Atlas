package engine

import "github.com/divijg19/Atlas/internal/index"

// RankingStrategy scores a candidate profile for search ranking. The concrete
// policy is owned by the evaluation layer; engine only defines the contract.
type RankingStrategy interface {
	Score(index.Profile) float64
}
