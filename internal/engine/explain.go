package engine

import (
	"fmt"
	"strings"

	"github.com/divijg19/Atlas/internal/index"
	"github.com/divijg19/Atlas/internal/signals"
)

func Explain(p index.Profile, q Query) []string {
	if len(q.Conditions) == 0 {
		return []string{"Ranked by overall signal strength"}
	}

	if p.Facts != nil {
		return explainWithEvidence(p, q)
	}

	return explainFromConditions(p, q)
}

func explainWithEvidence(p index.Profile, q Query) []string {
	sig := signals.FromMap(p.Signals)

	allEvidence := signals.GenerateEvidence(*p.Facts, sig)

	matched := map[string]struct{}{}
	for _, c := range q.Conditions {
		signal := strings.ToLower(strings.TrimSpace(c.Signal))
		if _, exists := matched[signal]; exists {
			continue
		}
		if Match(p, c) {
			matched[signal] = struct{}{}
		}
	}

	var reasons []string
	for _, group := range allEvidence {
		if _, ok := matched[group.Signal]; !ok {
			continue
		}
		for _, item := range group.Items {
			reasons = append(reasons, item.Description)
		}
	}

	if len(reasons) == 0 {
		return []string{"Ranked by overall signal strength"}
	}

	return reasons
}

func explainFromConditions(p index.Profile, q Query) []string {
	reasons := make([]string, 0, len(q.Conditions))
	seenSignals := map[string]struct{}{}

	for _, condition := range q.Conditions {
		if !Match(p, condition) {
			continue
		}

		signal := strings.ToLower(strings.TrimSpace(condition.Signal))
		reason := reasonForSignal(signal)
		if reason == "" {
			continue
		}
		if _, exists := seenSignals[signal]; exists {
			continue
		}

		reason = fmt.Sprintf("%s (%.2f %s %.2f)", reason, p.Signals[signal], strings.TrimSpace(condition.Operator), condition.Value)

		reasons = append(reasons, reason)
		seenSignals[signal] = struct{}{}
	}

	return reasons
}

func reasonForSignal(signal string) string {
	switch strings.ToLower(strings.TrimSpace(signal)) {
	case signals.SignalConsistency:
		return "High consistency"
	case signals.SignalOwnership:
		return "Strong ownership"
	case signals.SignalDepth:
		return "Deep project work"
	case signals.SignalActivity:
		return "Recently active"
	default:
		return ""
	}
}
