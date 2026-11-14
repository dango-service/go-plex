package plex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultPinsURL is the standard plex.tv endpoint for Plex PIN OAuth authentication.
	DefaultPinsURL = "https://plex.tv/api/v2/pins"
)

// AuthService handles Plex PIN OAuth authentication with plex.tv.
type AuthService struct {
	client *Client
}

// PIN represents a Plex OAuth PIN resource used to authenticate a client with plex.tv.
type PIN struct {
	ID               int64     `json:"id"`
	Code             string    `json:"code"`
	Product          string    `json:"product"`
	Trusted          bool      `json:"trusted"`
	ClientIdentifier string    `json:"clientIdentifier"`
	ExpiresIn        int       `json:"expiresIn"`
	CreatedAt        time.Time `json:"createdAt"`
	ExpiresAt        time.Time `json:"expiresAt"`
	AuthToken        string    `json:"authToken"`
	NewRegistration  *bool     `json:"newRegistration,omitempty"`
}

// IsClaimed returns true if the user has authenticated and claimed the PIN, providing an auth token.
func (p *PIN) IsClaimed() bool {
	return strings.TrimSpace(p.AuthToken) != ""
}

// IsExpired reports whether the PIN has expired.
func (p *PIN) IsExpired() bool {
	if p.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(p.ExpiresAt)
}

// AuthURL returns the standard Plex web app authorization URL that the user must visit to authorize the PIN.
func (p *PIN) AuthURL() string {
	return fmt.Sprintf("https://app.plex.tv/auth#?clientID=%s&code=%s",
		url.QueryEscape(p.ClientIdentifier),
		url.QueryEscape(p.Code),
	)
}

// CreatePIN creates a new Plex OAuth PIN on plex.tv.
// The user visits the returned AuthURL to authenticate their Plex account.
func (s *AuthService) CreatePIN(ctx context.Context) (*PIN, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	targetURL := s.client.pinsURL
	if targetURL == "" {
		targetURL = DefaultPinsURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("plex: create pin request: %w", err)
	}

	// Request a strong alphanumeric PIN code
	q := req.URL.Query()
	q.Set("strong", "true")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Product", s.client.product)
	req.Header.Set("X-Plex-Version", s.client.version)
	req.Header.Set("X-Plex-Client-Identifier", s.client.clientID)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("plex: execute create pin request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("plex: create pin returned HTTP %d: %w", resp.StatusCode, ErrBadRequest)
	}

	var pin PIN
	limitedReader := io.LimitReader(resp.Body, s.client.maxBodyBytes)
	if err := json.NewDecoder(limitedReader).Decode(&pin); err != nil {
		return nil, fmt.Errorf("plex: decode create pin response: %w", err)
	}

	return &pin, nil
}

// CheckPIN polls plex.tv for the authorization status of a previously generated PIN.
// If the user has claimed the PIN in the Plex web UI, the returned PIN will have AuthToken populated.
func (s *AuthService) CheckPIN(ctx context.Context, pinID int64) (*PIN, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if pinID <= 0 {
		return nil, fmt.Errorf("plex: invalid pin id %d: %w", pinID, ErrBadRequest)
	}

	targetURL := s.client.pinsURL
	if targetURL == "" {
		targetURL = DefaultPinsURL
	}

	fullURL := fmt.Sprintf("%s/%d", strings.TrimRight(targetURL, "/"), pinID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("plex: create check pin request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Client-Identifier", s.client.clientID)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("plex: execute check pin request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("plex: pin %d not found or expired: %w", pinID, ErrNotFound)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("plex: check pin returned HTTP %d", resp.StatusCode)
	}

	var pin PIN
	limitedReader := io.LimitReader(resp.Body, s.client.maxBodyBytes)
	if err := json.NewDecoder(limitedReader).Decode(&pin); err != nil {
		return nil, fmt.Errorf("plex: decode check pin response: %w", err)
	}

	return &pin, nil
}
