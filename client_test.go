package plex

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("http://localhost:32400/", "test-token")

	if c.baseURL != "http://localhost:32400" {
		t.Errorf("expected trimmed baseURL, got %q", c.baseURL)
	}
	if c.token != "test-token" {
		t.Errorf("expected token 'test-token', got %q", c.token)
	}
	if c.clientID != DefaultClientIdentifier {
		t.Errorf("expected clientID %q, got %q", DefaultClientIdentifier, c.clientID)
	}
	if c.product != DefaultProduct || c.version != DefaultVersion {
		t.Errorf("expected default product/version, got %s/%s", c.product, c.version)
	}
	if c.maxBodyBytes != DefaultMaxBodyBytes {
		t.Errorf("expected maxBodyBytes %d, got %d", DefaultMaxBodyBytes, c.maxBodyBytes)
	}
	if c.Library == nil || c.Status == nil || c.Discovery == nil {
		t.Fatal("expected services to be initialized")
	}
}

func TestNewClient_Options(t *testing.T) {
	customHTTP := &http.Client{Timeout: 5 * time.Second}
	c := NewClient("http://192.168.1.50:32400", "my-token",
		WithClientIdentifier("custom-client-id"),
		WithProduct("my-app"),
		WithVersion("2.1.0"),
		WithMaxBodyBytes(2*1024*1024),
		WithHTTPClient(customHTTP),
	)

	if c.clientID != "custom-client-id" {
		t.Errorf("expected custom clientID, got %q", c.clientID)
	}
	if c.product != "my-app" || c.version != "2.1.0" {
		t.Errorf("expected custom product/version, got %s/%s", c.product, c.version)
	}
	if c.maxBodyBytes != 2097152 {
		t.Errorf("expected 2MB maxBodyBytes, got %d", c.maxBodyBytes)
	}
	if c.httpClient != customHTTP {
		t.Fatal("expected custom http client to be set")
	}
}
