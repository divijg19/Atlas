package signals

// Canonical signal names. These are the single source of truth for the string
// keys used to identify the four portfolio indicators across the codebase
// (storage maps, query engine, ranking weights, evidence, presets).
const (
	SignalOwnership   = "ownership"
	SignalConsistency = "consistency"
	SignalDepth       = "depth"
	SignalActivity    = "activity"
)

// SignalsToMap converts a Signals struct to a map[string]float64
// suitable for storage in Profile and use by the query engine.
func SignalsToMap(s Signals) map[string]float64 {
	return map[string]float64{
		SignalOwnership:   clamp01(s.Ownership),
		SignalConsistency: clamp01(s.Consistency),
		SignalDepth:       clamp01(s.Depth),
		SignalActivity:    clamp01(s.Activity),
	}
}

// FromMap converts a map[string]float64 (as stored on a Profile) back into a
// Signals struct. It is the inverse of SignalsToMap and the single canonical
// owner of the signal-name → field mapping, eliminating duplicated literal
// key lookups elsewhere in the codebase.
func FromMap(m map[string]float64) Signals {
	return Signals{
		Ownership:   m[SignalOwnership],
		Consistency: m[SignalConsistency],
		Depth:       m[SignalDepth],
		Activity:    m[SignalActivity],
	}
}

func clamp01(value float64) float64 {
	if value < minSignalValue {
		return minSignalValue
	}
	if value > maxSignalValue {
		return maxSignalValue
	}

	return value
}
