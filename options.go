package plex

import (
	"net/http"
	"strings"
)

// ClientOption configures a Client instance.
type ClientOption func(*Client)

// WithHTTPClient overrides the default HTTP client used for requests.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithClientIdentifier sets the X-Plex-Client-Identifier header.
func WithClientIdentifier(clientID string) ClientOption {
	return func(c *Client) {
		trimmed := strings.TrimSpace(clientID)
		if trimmed != "" {
			c.clientID = trimmed
		}
	}
}

// WithProduct sets the X-Plex-Product header.
func WithProduct(product string) ClientOption {
	return func(c *Client) {
		trimmed := strings.TrimSpace(product)
		if trimmed != "" {
			c.product = trimmed
		}
	}
}

// WithVersion sets the X-Plex-Version header.
func WithVersion(version string) ClientOption {
	return func(c *Client) {
		trimmed := strings.TrimSpace(version)
		if trimmed != "" {
			c.version = trimmed
		}
	}
}

// WithMaxBodyBytes sets the maximum allowed response body size (default: 10 MiB).
func WithMaxBodyBytes(maxBytes int64) ClientOption {
	return func(c *Client) {
		if maxBytes > 0 {
			c.maxBodyBytes = maxBytes
		}
	}
}

// WithDiscoveryURL overrides the plex.tv resources discovery endpoint.
func WithDiscoveryURL(rawURL string) ClientOption {
	return func(c *Client) {
		trimmed := strings.TrimSpace(rawURL)
		if trimmed != "" {
			c.discoveryURL = trimmed
		}
	}
}

// WithPinsURL overrides the plex.tv PIN authentication endpoint (used for testing).
func WithPinsURL(rawURL string) ClientOption {
	return func(c *Client) {
		trimmed := strings.TrimSpace(rawURL)
		if trimmed != "" {
			c.pinsURL = trimmed
		}
	}
}
