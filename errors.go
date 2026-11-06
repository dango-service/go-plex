package plex

import (
	"errors"
	"fmt"
)

var (
	// ErrUnauthorized is returned when authentication is missing or invalid (HTTP 401).
	ErrUnauthorized = errors.New("plex: unauthorized (401)")

	// ErrNotFound is returned when the requested resource does not exist (HTTP 404).
	ErrNotFound = errors.New("plex: resource not found (404)")

	// ErrBadRequest is returned when request parameters are invalid (HTTP 400).
	ErrBadRequest = errors.New("plex: bad request (400)")

	// ErrIgnoredEvent is returned when a webhook event is skipped by filter criteria.
	ErrIgnoredEvent = errors.New("plex: webhook event ignored")
)

// ErrorResponse represents an HTTP error returned by Plex Media Server.
type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Message    string `json:"error,omitempty"`
	Err        error  `json:"-"`
}

func (e *ErrorResponse) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("plex: HTTP %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("plex: HTTP %d", e.StatusCode)
}

// Unwrap returns the underlying sentinel error for errors.Is and errors.As matching.
func (e *ErrorResponse) Unwrap() error {
	return e.Err
}

// IsUnauthorized checks if an error corresponds to HTTP 401.
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsNotFound checks if an error corresponds to HTTP 404.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
