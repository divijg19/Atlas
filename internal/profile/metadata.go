// Package profile defines the user metadata model owned by the domain layer.
//
// It holds normalized account observations (UserMetadata) produced by
// acquisition/normalization. It contains no networking and no scoring logic.
// See docs/INTELLIGENCE.md.
package profile

import (
	"time"
)

type UserMetadata struct {
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	Location  string    `json:"location"`
	Company   string    `json:"company"`
	Followers int       `json:"followers"`
	Following int       `json:"following"`
	CreatedAt time.Time `json:"created_at"`
}
