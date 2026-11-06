package plex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// DefaultDiscoveryURL is the standard plex.tv endpoint for querying user server resources.
	DefaultDiscoveryURL = "https://plex.tv/api/v2/resources?includeHttps=1&includeRelay=1"
)

// ServerConnection describes a network connection to a Plex Media Server.
type ServerConnection struct {
	URI     string `json:"uri"`
	Address string `json:"address"`
	Port    int    `json:"port"`
	Local   bool   `json:"local"`
	Relay   bool   `json:"relay"`
}

// ServerDevice describes a Plex device discovered via plex.tv.
type ServerDevice struct {
	Name             string             `json:"name"`
	Product          string             `json:"product"`
	Provides         string             `json:"provides"`
	ClientIdentifier string             `json:"clientIdentifier"`
	Connections      []ServerConnection `json:"connections"`
}

// BestConnection selects the best available URI (prefers local direct connection over relay).
func (d *ServerDevice) BestConnection() string {
	// 1. Prefer local non-relay connection
	for _, c := range d.Connections {
		if c.Local && !c.Relay && c.URI != "" {
			return strings.TrimRight(c.URI, "/")
		}
	}
	// 2. Prefer remote direct non-relay connection
	for _, c := range d.Connections {
		if !c.Relay && c.URI != "" {
			return strings.TrimRight(c.URI, "/")
		}
	}
	// 3. Fallback to relay or first available
	if len(d.Connections) > 0 && d.Connections[0].URI != "" {
		return strings.TrimRight(d.Connections[0].URI, "/")
	}
	return ""
}

// DiscoverServers queries plex.tv to discover all available server devices associated with a token.
func (s *DiscoveryService) DiscoverServers(ctx context.Context, plexToken string) ([]ServerDevice, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	targetURL := s.client.discoveryURL
	if targetURL == "" {
		targetURL = DefaultDiscoveryURL
	}

	token := strings.TrimSpace(plexToken)
	if token == "" {
		token = s.client.token
	}
	if token == "" {
		return nil, fmt.Errorf("plex: plex token is required for server discovery")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("plex: create discovery request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Token", token)
	req.Header.Set("X-Plex-Client-Identifier", s.client.clientID)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("plex: execute discovery request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("plex: discovery returned HTTP %d", resp.StatusCode)
	}

	var devices []ServerDevice
	limitedReader := io.LimitReader(resp.Body, s.client.maxBodyBytes)
	if err := json.NewDecoder(limitedReader).Decode(&devices); err != nil {
		return nil, fmt.Errorf("plex: decode discovery response: %w", err)
	}

	return devices, nil
}

// EnsureActiveURL verifies if the current baseURL is reachable; if not, discovers active connection from plex.tv.
func (s *DiscoveryService) EnsureActiveURL(ctx context.Context) string {
	if s.client.baseURL != "" && s.client.token != "" {
		testReq, err := http.NewRequestWithContext(ctx, http.MethodGet, s.client.baseURL+"/library/sections", nil)
		if err == nil {
			testReq.Header.Set("X-Plex-Token", s.client.token)
			testReq.Header.Set("Accept", "application/json")
			resp, err := s.client.httpClient.Do(testReq)
			if err == nil {
				_ = resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					return s.client.baseURL
				}
			}
		}
	}

	// Discover live URI from plex.tv
	devices, err := s.DiscoverServers(ctx, s.client.token)
	if err != nil {
		return s.client.baseURL
	}

	for _, dev := range devices {
		if strings.Contains(dev.Provides, "server") {
			for _, conn := range dev.Connections {
				candidate := strings.TrimRight(conn.URI, "/")
				if candidate == "" {
					continue
				}
				cReq, err := http.NewRequestWithContext(ctx, http.MethodGet, candidate+"/library/sections", nil)
				if err == nil {
					cReq.Header.Set("X-Plex-Token", s.client.token)
					cReq.Header.Set("Accept", "application/json")
					cResp, cErr := s.client.httpClient.Do(cReq)
					if cErr == nil {
						_ = cResp.Body.Close()
						if cResp.StatusCode == http.StatusOK {
							s.client.baseURL = candidate
							return candidate
						}
					}
				}
			}
		}
	}

	return s.client.baseURL
}
