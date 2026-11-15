package plex

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestLibraryService_GetSections(t *testing.T) {
	fixture, err := os.ReadFile("testdata/library_sections_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/library/sections" {
			t.Errorf("expected path '/library/sections', got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	sections, resp, err := client.Library.GetSections(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(sections))
	}
	if sections[0].Key != "1" || sections[0].Title != "Anime" || sections[0].Type != "show" {
		t.Errorf("unexpected section 0: %+v", sections[0])
	}
}

func TestLibraryService_GetSectionItems(t *testing.T) {
	fixture, err := os.ReadFile("testdata/library_section_items_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/library/sections/1/all" {
			t.Errorf("expected path '/library/sections/1/all', got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	items, resp, err := client.Library.GetSectionItems(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Title != "Sousou no Frieren" || items[0].TVDB() != "404080" {
		t.Errorf("unexpected item: %+v", items[0])
	}
}

func TestLibraryService_GetSectionItems_Validation(t *testing.T) {
	client := NewClient("http://localhost:32400", "test-token")
	if _, _, err := client.Library.GetSectionItems(context.Background(), "   "); err == nil {
		t.Fatal("expected error on empty section key, got nil")
	}
}

func TestLibraryService_GetRecentlyWatchedItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/library/sections/1/all" || r.URL.Query().Get("type") != "4" ||
			r.URL.Query().Get("lastViewedAt>>") != "2026-11-15" || r.URL.Query().Get("includeGuids") != "1" {
			t.Errorf("unexpected recent episode query: %s", r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"MediaContainer":{"Metadata":[{"ratingKey":"5","type":"episode","viewCount":1,"lastViewedAt":1794758400}]}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	since := time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)
	items, _, err := client.Library.GetRecentlyWatchedItems(context.Background(), "1", "show", since)
	if err != nil || len(items) != 1 || items[0].RatingKey.String() != "5" {
		t.Fatalf("unexpected recently watched episodes: items=%+v err=%v", items, err)
	}
	if _, _, err := client.Library.GetRecentlyWatchedItems(context.Background(), "1", "artist", since); err == nil {
		t.Fatal("unsupported section type should fail")
	}
}

func TestLibraryService_GetMetadata(t *testing.T) {
	fixture, err := os.ReadFile("testdata/library_metadata_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/library/metadata/5001" {
			t.Errorf("expected path '/library/metadata/5001', got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	item, resp, err := client.Library.GetMetadata(context.Background(), "5001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if item.RatingKey.String() != "5001" || item.Title != "The Journey's End" {
		t.Errorf("unexpected metadata: %+v", item)
	}
	if item.Index != 1 || item.ParentIndex != 1 || item.GrandparentTitle != "Sousou no Frieren" {
		t.Errorf("unexpected episode indices: %+v", item)
	}
}

func TestLibraryService_GetMetadata_Validation(t *testing.T) {
	client := NewClient("http://localhost:32400", "test-token")
	if _, _, err := client.Library.GetMetadata(context.Background(), ""); err == nil {
		t.Fatal("expected error on empty rating key, got nil")
	}
}

func TestLibraryService_GetAllLeaves(t *testing.T) {
	fixture, err := os.ReadFile("testdata/library_all_leaves_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/library/metadata/1001/allLeaves" {
			t.Errorf("expected path '/library/metadata/1001/allLeaves', got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	leaves, resp, err := client.Library.GetAllLeaves(context.Background(), "1001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if len(leaves) != 2 {
		t.Fatalf("expected 2 leaves, got %d", len(leaves))
	}
	if leaves[0].Index != 1 || leaves[1].Index != 2 {
		t.Errorf("unexpected leaves indices: %+v", leaves)
	}
}

func TestLibraryService_GetAllLeaves_Validation(t *testing.T) {
	client := NewClient("http://localhost:32400", "test-token")
	if _, _, err := client.Library.GetAllLeaves(context.Background(), "  "); err == nil {
		t.Fatal("expected error on empty rating key, got nil")
	}
}
