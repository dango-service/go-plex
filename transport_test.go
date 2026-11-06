package plex

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTransport_HeadersAndSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Plex-Token") != "secret-token" {
			t.Errorf("expected X-Plex-Token 'secret-token', got %q", r.Header.Get("X-Plex-Token"))
		}
		if r.Header.Get("X-Plex-Client-Identifier") != DefaultClientIdentifier {
			t.Errorf("expected X-Plex-Client-Identifier, got %q", r.Header.Get("X-Plex-Client-Identifier"))
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept application/json, got %q", r.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "secret-token")
	var out struct {
		Status string `json:"status"`
	}

	resp, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if out.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", out.Status)
	}
}

func TestTransport_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "invalid-token")
	_, err := client.doRequest(context.Background(), http.MethodGet, "/secure", nil, nil)

	if !IsUnauthorized(err) {
		t.Fatalf("expected IsUnauthorized(err) to be true, got %v", err)
	}
}

func TestTransport_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "valid-token")
	_, err := client.doRequest(context.Background(), http.MethodGet, "/missing", nil, nil)

	if !IsNotFound(err) {
		t.Fatalf("expected IsNotFound(err) to be true, got %v", err)
	}
}

func TestTransport_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "valid-token")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.doRequest(ctx, http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error on cancelled context, got nil")
	}
}
