package plex

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Section represents a Plex library section (e.g. Anime, TV Shows, Movies).
type Section struct {
	Key       string `json:"key"`
	Title     string `json:"title"`
	Type      string `json:"type"` // "show", "movie", "artist", "photo"
	Language  string `json:"language,omitempty"`
	Agent     string `json:"agent,omitempty"`
	Scanner   string `json:"scanner,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
}

// GetSections retrieves all library sections configured on the Plex Media Server.
func (s *LibraryService) GetSections(ctx context.Context) ([]Section, *Response, error) {
	var mc MediaContainer[Section]
	resp, err := s.client.doRequest(ctx, http.MethodGet, "/library/sections", nil, &mc)
	if err != nil {
		return nil, resp, fmt.Errorf("plex: get sections: %w", err)
	}

	return mc.Items(), resp, nil
}

// GetSectionItems retrieves all items in the specified library section key.
func (s *LibraryService) GetSectionItems(ctx context.Context, sectionKey string) ([]MetadataItem, *Response, error) {
	trimmed := strings.TrimSpace(sectionKey)
	if trimmed == "" {
		return nil, nil, fmt.Errorf("plex: section key cannot be empty")
	}

	path := fmt.Sprintf("/library/sections/%s/all", trimmed)
	var mc MediaContainer[MetadataItem]
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &mc)
	if err != nil {
		return nil, resp, fmt.Errorf("plex: get section items: %w", err)
	}

	return mc.Items(), resp, nil
}

// GetMetadata retrieves full metadata for an item given its ratingKey.
func (s *LibraryService) GetMetadata(ctx context.Context, ratingKey string) (*MetadataItem, *Response, error) {
	trimmed := strings.TrimSpace(ratingKey)
	if trimmed == "" {
		return nil, nil, fmt.Errorf("plex: rating key cannot be empty")
	}

	path := fmt.Sprintf("/library/metadata/%s", trimmed)
	var mc MediaContainer[MetadataItem]
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &mc)
	if err != nil {
		return nil, resp, fmt.Errorf("plex: get metadata: %w", err)
	}

	items := mc.Items()
	if len(items) == 0 {
		return nil, resp, fmt.Errorf("plex: metadata not found for rating key %s: %w", trimmed, ErrNotFound)
	}

	return &items[0], resp, nil
}

// GetAllLeaves retrieves all episodes belonging to a show or season given its ratingKey.
func (s *LibraryService) GetAllLeaves(ctx context.Context, ratingKey string) ([]MetadataItem, *Response, error) {
	trimmed := strings.TrimSpace(ratingKey)
	if trimmed == "" {
		return nil, nil, fmt.Errorf("plex: rating key cannot be empty")
	}

	path := fmt.Sprintf("/library/metadata/%s/allLeaves", trimmed)
	var mc MediaContainer[MetadataItem]
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &mc)
	if err != nil {
		return nil, resp, fmt.Errorf("plex: get all leaves: %w", err)
	}

	return mc.Items(), resp, nil
}
