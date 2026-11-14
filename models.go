package plex

import (
	"encoding/json"
	"strconv"
	"strings"
)

// FlexibleID unmarshals IDs that may appear as numbers or strings in Plex JSON responses.
type FlexibleID string

func (f *FlexibleID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*f = FlexibleID(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(b, &n); err == nil {
		*f = FlexibleID(n.String())
		return nil
	}
	return nil
}

// String returns the string representation of the ID.
func (f FlexibleID) String() string {
	return string(f)
}

// Int returns the integer representation of the ID, or 0 if parsing fails.
func (f FlexibleID) Int() int {
	i, _ := strconv.Atoi(string(f))
	return i
}

// Guid represents an external identifier provided by Plex metadata agents.
type Guid struct {
	ID string `json:"id"`
}

// Scheme returns the prefix scheme of the GUID (e.g. "tvdb", "tmdb", "imdb", "anidb").
func (g Guid) Scheme() string {
	parts := strings.SplitN(g.ID, "://", 2)
	if len(parts) == 2 {
		scheme := strings.ToLower(strings.TrimSpace(parts[0]))
		scheme = strings.TrimPrefix(scheme, "com.plexapp.agents.")
		switch scheme {
		case "thetvdb":
			return "tvdb"
		case "themoviedb":
			return "tmdb"
		default:
			return scheme
		}
	}
	return ""
}

// Value returns the identifier value stripped of the scheme prefix.
func (g Guid) Value() string {
	parts := strings.SplitN(g.ID, "://", 2)
	if len(parts) == 2 {
		value := parts[1]
		if end := strings.IndexAny(value, "?#"); end >= 0 {
			value = value[:end]
		}
		return strings.TrimSpace(value)
	}
	return g.ID
}

// TagItem represents a generic tag attribute (Genre, Collection, Tag).
type TagItem struct {
	Tag string `json:"tag"`
}

// MetadataItem represents a media item in Plex Media Server.
type MetadataItem struct {
	RatingKey                FlexibleID `json:"ratingKey"`
	Key                      string     `json:"key"`
	Guid                     string     `json:"guid"`
	Type                     string     `json:"type"` // "episode", "movie", "show", "season"
	Title                    string     `json:"title"`
	OriginalTitle            string     `json:"originalTitle,omitempty"`
	Year                     int        `json:"year,omitempty"`
	Thumb                    string     `json:"thumb,omitempty"`
	Art                      string     `json:"art,omitempty"`
	Duration                 int64      `json:"duration,omitempty"`
	ViewOffset               int64      `json:"viewOffset,omitempty"`
	ViewCount                int        `json:"viewCount,omitempty"`
	LastViewedAt             int64      `json:"lastViewedAt,omitempty"`
	UserRating               *float64   `json:"userRating,omitempty"`
	Index                    int        `json:"index,omitempty"`       // Episode number
	ParentIndex              int        `json:"parentIndex,omitempty"` // Season number
	ParentTitle              string     `json:"parentTitle,omitempty"`
	ParentRatingKey          FlexibleID `json:"parentRatingKey,omitempty"`
	GrandparentTitle         string     `json:"grandparentTitle,omitempty"`
	GrandparentOriginalTitle string     `json:"grandparentOriginalTitle,omitempty"`
	GrandparentRatingKey     FlexibleID `json:"grandparentRatingKey,omitempty"`
	ParentGuid               string     `json:"parentGuid,omitempty"`
	GrandparentGuid          string     `json:"grandparentGuid,omitempty"`
	GrandparentThumb         string     `json:"grandparentThumb,omitempty"`
	LibrarySectionID         FlexibleID `json:"librarySectionID,omitempty"`
	LibrarySectionTitle      string     `json:"librarySectionTitle,omitempty"`
	Guids                    []Guid     `json:"Guid,omitempty"`
	Genre                    []TagItem  `json:"Genre,omitempty"`
	Tag                      []TagItem  `json:"Tag,omitempty"`
	Collection               []TagItem  `json:"Collection,omitempty"`
}

func (m *MetadataItem) externalID(scheme string) string {
	if m == nil {
		return ""
	}
	guids := append([]Guid(nil), m.Guids...)
	guids = append(guids, Guid{ID: m.Guid}, Guid{ID: m.ParentGuid}, Guid{ID: m.GrandparentGuid})
	for _, g := range guids {
		if g.Scheme() == scheme {
			return g.Value()
		}
	}
	return ""
}

// TVDB returns the TVDB ID from an item's GUIDs or its ancestor GUID fields.
func (m *MetadataItem) TVDB() string {
	return m.externalID("tvdb")
}

// TMDB returns the TMDB ID if present in Guids.
func (m *MetadataItem) TMDB() string {
	return m.externalID("tmdb")
}

// IMDb returns the IMDb ID if present in Guids.
func (m *MetadataItem) IMDb() string {
	return m.externalID("imdb")
}

// AniDB returns the AniDB ID if present in Guids.
func (m *MetadataItem) AniDB() string {
	return m.externalID("anidb")
}

// GenreNames extracts string genre names.
func (m *MetadataItem) GenreNames() []string {
	names := make([]string, 0, len(m.Genre))
	for _, g := range m.Genre {
		if g.Tag != "" {
			names = append(names, g.Tag)
		}
	}
	return names
}

// TagNames extracts string tag names.
func (m *MetadataItem) TagNames() []string {
	names := make([]string, 0, len(m.Tag))
	for _, t := range m.Tag {
		if t.Tag != "" {
			names = append(names, t.Tag)
		}
	}
	return names
}

// CollectionNames extracts string collection names.
func (m *MetadataItem) CollectionNames() []string {
	names := make([]string, 0, len(m.Collection))
	for _, c := range m.Collection {
		if c.Tag != "" {
			names = append(names, c.Tag)
		}
	}
	return names
}

// MediaContainer wraps standard Plex JSON response structures.
type MediaContainer[T any] struct {
	MediaContainer struct {
		Size      int `json:"size"`
		Metadata  []T `json:"Metadata,omitempty"`
		Directory []T `json:"Directory,omitempty"`
	} `json:"MediaContainer"`
}

// Items returns either Metadata or Directory slices from the container.
func (mc *MediaContainer[T]) Items() []T {
	if len(mc.MediaContainer.Metadata) > 0 {
		return mc.MediaContainer.Metadata
	}
	return mc.MediaContainer.Directory
}
