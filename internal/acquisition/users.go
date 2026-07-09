package acquisition

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserDTO mirrors a user object from the GitHub REST API. Field names and types
// follow GitHub's schema; created_at remains a string and is parsed during
// normalization.
type UserDTO struct {
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Company   string `json:"company"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
	CreatedAt string `json:"created_at"`
}

// FetchUser retrieves a user's public metadata from the GitHub REST API.
func (c *Client) FetchUser(ctx context.Context, username string) (*UserDTO, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	resp, err := c.get(ctx, fmt.Sprintf("%s/users/%s", c.baseURL, username))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var ghErr struct {
			Message string `json:"message"`
		}

		message := resp.Status
		if err := json.NewDecoder(resp.Body).Decode(&ghErr); err == nil && ghErr.Message != "" {
			message = ghErr.Message
		}

		return nil, fmt.Errorf("GitHub API error (status %d): %s", resp.StatusCode, message)
	}

	var dto UserDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, fmt.Errorf("decode user metadata: %w", err)
	}

	return &dto, nil
}
