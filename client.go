package plex

import (
	"net/http"
	"strings"
	"time"
)

const (
	// DefaultProduct is sent in the X-Plex-Product header.
	DefaultProduct = "go-plex"

	// DefaultVersion is sent in the X-Plex-Version header.
	DefaultVersion = "0.1.0"

	// DefaultClientIdentifier is sent in the X-Plex-Client-Identifier header.
	DefaultClientIdentifier = "dango-plex-client"

	// DefaultMaxBodyBytes limits incoming JSON payloads to 10 MiB.
	DefaultMaxBodyBytes = 10 * 1024 * 1024
)

// Client is the Plex Media Server API client.
type Client struct {
	httpClient   *http.Client
	baseURL      string
	discoveryURL string
	token        string
	clientID     string
	product      string
	version      string
	maxBodyBytes int64

	Library   *LibraryService
	Status    *StatusService
	Discovery *DiscoveryService
}

// NewClient creates a new Plex Media Server client.
func NewClient(baseURL, token string, opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL:      strings.TrimRight(baseURL, "/"),
		discoveryURL: DefaultDiscoveryURL,
		token:        strings.TrimSpace(token),
		clientID:     DefaultClientIdentifier,
		product:      DefaultProduct,
		version:      DefaultVersion,
		maxBodyBytes: DefaultMaxBodyBytes,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}

	c.Library = &LibraryService{client: c}
	c.Status = &StatusService{client: c}
	c.Discovery = &DiscoveryService{client: c}

	return c
}

// LibraryService handles library sections, metadata, and episode leaves.
type LibraryService struct {
	client *Client
}

// StatusService handles active playback sessions and watch history.
type StatusService struct {
	client *Client
}

// DiscoveryService discovers active Plex server connections via plex.tv.
type DiscoveryService struct {
	client *Client
}
