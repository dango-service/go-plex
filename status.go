package plex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// SessionUser identifies the user account playing a session.
type SessionUser struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// SessionPlayer identifies the player client playing a session.
type SessionPlayer struct {
	Title string `json:"title"`
	State string `json:"state"` // "playing", "paused", "buffering"
}

// SessionItem represents an active streaming session on Plex Media Server.
type SessionItem struct {
	MetadataItem
	User   *SessionUser   `json:"User,omitempty"`
	Player *SessionPlayer `json:"Player,omitempty"`
}

// HistoryItem represents an entry from the server watch history.
type HistoryItem struct {
	HistoryKey   FlexibleID `json:"historyKey,omitempty"`
	RatingKey    FlexibleID `json:"ratingKey"`
	Key          string     `json:"key"`
	Title        string     `json:"title"`
	Type         string     `json:"type"` // "episode", "movie"
	ViewedAt     int64      `json:"viewedAt"`
	AccountID    int        `json:"accountID,omitempty"`
	AccountTitle string     `json:"accountTitle,omitempty"`
}

// GetSessions retrieves all active streaming playback sessions.
func (s *StatusService) GetSessions(ctx context.Context) ([]SessionItem, *Response, error) {
	var mc MediaContainer[SessionItem]
	resp, err := s.client.doRequest(ctx, http.MethodGet, "/status/sessions", nil, &mc)
	if err != nil {
		return nil, resp, fmt.Errorf("plex: get sessions: %w", err)
	}

	return mc.Items(), resp, nil
}

// GetHistory retrieves the server watch history. If minViewedAt is > 0, only records watched since then are returned.
func (s *StatusService) GetHistory(ctx context.Context, minViewedAt int64) ([]HistoryItem, *Response, error) {
	path := "/status/sessions/history/all"
	if minViewedAt > 0 {
		v := url.Values{}
		v.Set("viewedAt>=", strconv.FormatInt(minViewedAt, 10))
		path += "?" + v.Encode()
	}

	var mc MediaContainer[HistoryItem]
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &mc)
	if err != nil {
		return nil, resp, fmt.Errorf("plex: get history: %w", err)
	}

	return mc.Items(), resp, nil
}
