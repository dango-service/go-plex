package plex

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestStatusService_GetSessions(t *testing.T) {
	fixture, err := os.ReadFile("testdata/status_sessions_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/status/sessions" {
			t.Errorf("expected path '/status/sessions', got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	sessions, resp, err := client.Status.GetSessions(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Title != "The Journey's End" || sessions[0].User == nil || sessions[0].User.Title != "hicham" {
		t.Errorf("unexpected session data: %+v", sessions[0])
	}
	if sessions[0].Player == nil || sessions[0].Player.State != "playing" {
		t.Errorf("unexpected player data: %+v", sessions[0].Player)
	}
}

func TestStatusService_GetHistory(t *testing.T) {
	fixture, err := os.ReadFile("testdata/status_history_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/status/sessions/history/all" {
			t.Errorf("expected path '/status/sessions/history/all', got %q", r.URL.Path)
		}
		if r.URL.Query().Get("viewedAt>=") != "1699999999" {
			t.Errorf("expected viewedAt>= filter, got %q", r.URL.Query().Get("viewedAt>="))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	history, resp, err := client.Status.GetHistory(context.Background(), 1699999999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if len(history) != 1 {
		t.Fatalf("expected 1 history item, got %d", len(history))
	}
	if history[0].RatingKey.String() != "5001" || history[0].AccountTitle != "hicham" {
		t.Errorf("unexpected history item: %+v", history[0])
	}
}
