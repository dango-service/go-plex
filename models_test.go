package plex

import (
	"encoding/json"
	"testing"
)

func TestFlexibleID_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		expected string
		asInt    int
	}{
		{"numeric id", `{"id": 12345}`, "12345", 12345},
		{"string id", `{"id": "67890"}`, "67890", 67890},
		{"zero numeric", `{"id": 0}`, "0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wrapper struct {
				ID FlexibleID `json:"id"`
			}
			if err := json.Unmarshal([]byte(tt.jsonStr), &wrapper); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}
			if wrapper.ID.String() != tt.expected {
				t.Errorf("expected string %q, got %q", tt.expected, wrapper.ID.String())
			}
			if wrapper.ID.Int() != tt.asInt {
				t.Errorf("expected int %d, got %d", tt.asInt, wrapper.ID.Int())
			}
		})
	}
}

func TestGuid_SchemeAndValue(t *testing.T) {
	g1 := Guid{ID: "tvdb://12345"}
	if g1.Scheme() != "tvdb" || g1.Value() != "12345" {
		t.Errorf("unexpected guid 1 parsing: %s / %s", g1.Scheme(), g1.Value())
	}

	g2 := Guid{ID: "imdb://tt0903747"}
	if g2.Scheme() != "imdb" || g2.Value() != "tt0903747" {
		t.Errorf("unexpected guid 2 parsing: %s / %s", g2.Scheme(), g2.Value())
	}

	g3 := Guid{ID: "rawid"}
	if g3.Scheme() != "" || g3.Value() != "rawid" {
		t.Errorf("unexpected guid 3 parsing: %s / %s", g3.Scheme(), g3.Value())
	}
}

func TestMetadataItem_Helpers(t *testing.T) {
	item := MetadataItem{
		RatingKey: "100",
		Title:     "Sousou no Frieren",
		Guids: []Guid{
			{ID: "tvdb://404080"},
			{ID: "tmdb://209867"},
			{ID: "imdb://tt22064048"},
			{ID: "anidb://17617"},
		},
		Genre: []TagItem{
			{Tag: "Animation"},
			{Tag: "Adventure"},
			{Tag: "Drama"},
		},
		Tag: []TagItem{
			{Tag: "Anime"},
		},
		Collection: []TagItem{
			{Tag: "Top Rated Anime"},
		},
	}

	if item.TVDB() != "404080" {
		t.Errorf("expected TVDB 404080, got %q", item.TVDB())
	}
	if item.TMDB() != "209867" {
		t.Errorf("expected TMDB 209867, got %q", item.TMDB())
	}
	if item.IMDb() != "tt22064048" {
		t.Errorf("expected IMDb tt22064048, got %q", item.IMDb())
	}
	if item.AniDB() != "17617" {
		t.Errorf("expected AniDB 17617, got %q", item.AniDB())
	}

	genres := item.GenreNames()
	if len(genres) != 3 || genres[1] != "Adventure" {
		t.Errorf("unexpected genres: %v", genres)
	}

	tags := item.TagNames()
	if len(tags) != 1 || tags[0] != "Anime" {
		t.Errorf("unexpected tags: %v", tags)
	}

	colls := item.CollectionNames()
	if len(colls) != 1 || colls[0] != "Top Rated Anime" {
		t.Errorf("unexpected collections: %v", colls)
	}
}

func TestMediaContainer_Items(t *testing.T) {
	var mc MediaContainer[MetadataItem]
	jsonBytes := []byte(`{
		"MediaContainer": {
			"size": 1,
			"Metadata": [
				{"ratingKey": "1", "title": "Episode 1"}
			]
		}
	}`)
	if err := json.Unmarshal(jsonBytes, &mc); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	items := mc.Items()
	if len(items) != 1 || items[0].Title != "Episode 1" {
		t.Errorf("unexpected items: %+v", items)
	}
}
