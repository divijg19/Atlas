package acquisition

// APIError represents a non-successful GitHub REST API response. It carries the
// HTTP status code so callers can map GitHub failures to their own responses.
type APIError struct {
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	return e.Message
}
