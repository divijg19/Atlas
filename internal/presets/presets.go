package presets

import (
	"fmt"

	"github.com/divijg19/Atlas/internal/signals"
)

type Query struct {
	Conditions []Condition
	Limit      int
}

type Condition struct {
	Signal   string
	Operator string
	Value    float64
}

func Preset(name string) (Query, error) {
	switch name {
	case "strong":
		return Query{
			Conditions: []Condition{
				{Signal: signals.SignalConsistency, Operator: ">=", Value: 0.7},
				{Signal: signals.SignalOwnership, Operator: ">=", Value: 0.6},
			},
		}, nil
	case "consistent":
		return Query{
			Conditions: []Condition{
				{Signal: signals.SignalConsistency, Operator: ">=", Value: 0.8},
			},
		}, nil
	case "deep":
		return Query{
			Conditions: []Condition{
				{Signal: signals.SignalDepth, Operator: ">=", Value: 0.7},
			},
		}, nil
	default:
		return Query{}, fmt.Errorf("unknown preset %q", name)
	}
}
