package plex

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDiscoveryService_DiscoverServers(t *testing.T) {
	fixture, err := os.ReadFile("testdata/discovery_resources_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Plex-Token") != "my-token" {
			t.Errorf("expected X-Plex-Token 'my-token', got %q", r.Header.Get("X-Plex-Token"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("http://localhost:32400", "my-token", WithDiscoveryURL(server.URL))
	devices, err := client.Discovery.DiscoverServers(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected discovery error: %v", err)
	}

	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	dev := devices[0]
	if dev.Name != "Tower-PMS" || dev.Provides != "server" {
		t.Errorf("unexpected device: %+v", dev)
	}
	best := dev.BestConnection()
	if best != "https://192-168-1-100.local.plex.direct:32400" {
		t.Errorf("expected local non-relay connection, got %q", best)
	}
}

func TestDiscoveryService_EnsureActiveURL(t *testing.T) {
	liveServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer liveServer.Close()

	// Initial dead server
	client := NewClient(liveServer.URL, "token")
	active := client.Discovery.EnsureActiveURL(context.Background())

	if active != liveServer.URL {
		t.Errorf("expected active URL %q, got %q", liveServer.URL, active)
	}
}
