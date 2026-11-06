package plex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Standard Plex webhook event names.
const (
	EventMediaPlay     = "media.play"
	EventMediaPause    = "media.pause"
	EventMediaResume   = "media.resume"
	EventMediaStop     = "media.stop"
	EventMediaScrobble = "media.scrobble"
	EventMediaRate     = "media.rate"
)

// WebhookAccount identifies the account triggering a webhook event.
type WebhookAccount struct {
	ID    int    `json:"id"`
	Thumb string `json:"thumb,omitempty"`
	Title string `json:"title"`
}

// WebhookServer identifies the Plex server generating a webhook event.
type WebhookServer struct {
	Title string `json:"title"`
	UUID  string `json:"uuid"`
}

// WebhookPlayer identifies the player client playing media in a webhook event.
type WebhookPlayer struct {
	Local         bool   `json:"local"`
	PublicAddress string `json:"publicAddress,omitempty"`
	Title         string `json:"title"`
	UUID          string `json:"uuid"`
}

// WebhookPayload represents the deserialized JSON payload delivered in Plex webhooks.
type WebhookPayload struct {
	Event    string         `json:"event"`
	User     bool           `json:"user"`
	Owner    bool           `json:"owner"`
	Rating   *float64       `json:"rating,omitempty"`
	Account  WebhookAccount `json:"Account"`
	Server   WebhookServer  `json:"Server"`
	Player   WebhookPlayer  `json:"Player"`
	Metadata MetadataItem   `json:"Metadata"`
}

// IsScrobble returns true if the event represents a completed media play.
func (p *WebhookPayload) IsScrobble() bool {
	return p.Event == EventMediaScrobble
}

// IsPlay returns true if playback started.
func (p *WebhookPayload) IsPlay() bool {
	return p.Event == EventMediaPlay
}

// IsPause returns true if playback paused.
func (p *WebhookPayload) IsPause() bool {
	return p.Event == EventMediaPause
}

// IsStop returns true if playback stopped.
func (p *WebhookPayload) IsStop() bool {
	return p.Event == EventMediaStop
}

// IsRate returns true if an item rating was updated.
func (p *WebhookPayload) IsRate() bool {
	return p.Event == EventMediaRate
}

// ParseWebhook processes an incoming Plex multipart/form-data webhook request.
// It discards thumbnail images on the fly via io.Discard to ensure peak memory
// efficiency, and decodes the JSON payload with bounded streaming.
func ParseWebhook(r *http.Request) (*WebhookPayload, error) {
	if r == nil {
		return nil, fmt.Errorf("plex: http request cannot be nil")
	}

	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("plex: invalid webhook method %s: %w", r.Method, ErrBadRequest)
	}

	mr, err := r.MultipartReader()
	if err != nil {
		return nil, fmt.Errorf("plex: read multipart body: %w", err)
	}

	var payload WebhookPayload
	var payloadFound bool

	for {
		part, err := mr.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("plex: read multipart part: %w", err)
		}

		formName := part.FormName()
		if formName == "thumb" {
			// Zero-allocation: discard image bytes without buffering into memory
			_, _ = io.Copy(io.Discard, part)
			_ = part.Close()
			continue
		}

		if formName == "payload" {
			limitedReader := io.LimitReader(part, DefaultMaxBodyBytes)
			decoder := json.NewDecoder(limitedReader)
			if err := decoder.Decode(&payload); err != nil {
				_ = part.Close()
				return nil, fmt.Errorf("plex: decode webhook json: %w", err)
			}
			payloadFound = true
			_ = part.Close()
			continue
		}

		// Discard any unrecognized parts
		_, _ = io.Copy(io.Discard, part)
		_ = part.Close()
	}

	if !payloadFound {
		return nil, fmt.Errorf("plex: no payload part found in multipart request: %w", ErrBadRequest)
	}

	return &payload, nil
}
