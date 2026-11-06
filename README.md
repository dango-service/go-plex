# go-plex

[![Go Reference](https://pkg.go.dev/badge/github.com/dango-service/go-plex.svg)](https://pkg.go.dev/github.com/dango-service/go-plex)
[![CI](https://github.com/dango-service/go-plex/actions/workflows/ci.yml/badge.svg)](https://github.com/dango-service/go-plex/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A lightweight, standalone Go client and webhook processing library for [Plex Media Server](https://www.plex.tv).
Designed with **zero external dependencies** (Go standard library only), direct streaming deserialization bounded by memory limits, and strict efficiency for background sync daemons like [Dango](https://github.com/dango-service/dango).

---

## Features

- **Zero External Dependencies**: Built entirely upon the Go standard library (`net/http`, `mime/multipart`, `encoding/json`, `context`, `net/url`).
- **Dual Architecture**:
  - **Plex Webhooks**: Real-time multipart webhook parser with zero-allocation thumbnail streaming straight to `io.Discard`.
  - **Internal PMS HTTP API**: Direct communication with local or remote Plex Media Servers for polling, initial library scans, and catch-up synchronization.
- **Streaming JSON**: Direct streaming decoding with `json.NewDecoder` bounded by `io.LimitReader` (10 MiB ceiling) to prevent unbounded memory growth and memory exhaustion attacks.
- **Flexible ID Resilience**: Safely deserializes Plex IDs that appear inconsistently as either integers or strings (`ratingKey`, `librarySectionID`) across PMS versions.
- **External GUID Helpers**: First-class helpers to resolve external IDs directly from media records: `item.TVDB()`, `item.TMDB()`, `item.IMDb()`, and `item.AniDB()`.
- **Active Sessions & History**: Inspect real-time streaming playback sessions (`/status/sessions`) and server watch history (`/status/sessions/history/all`).
- **Server Discovery**: Optional discovery of live connection URIs via `plex.tv/api/v2/resources`.
- **100% Deterministic Testing**: Complete test suite using `net/http/httptest` with zero real-network requests.

---

## Installation

```bash
go get github.com/dango-service/go-plex
```

---

## Usage Examples

### 1. Client Initialization

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/dango-service/go-plex"
)

func main() {
    client := plex.NewClient("http://192.168.1.100:32400", "YOUR_PLEX_TOKEN")

    // Fetch all library sections
    sections, _, err := client.Library.GetSections(context.Background())
    if err != nil {
        log.Fatalf("failed to fetch sections: %v", err)
    }

    for _, sec := range sections {
        fmt.Printf("Section [%s]: %s (Type: %s)\n", sec.Key, sec.Title, sec.Type)
    }
}
```

### 2. Inspecting Library Items & Episode Leaves

```go
// Fetch all shows or movies in a section
items, _, err := client.Library.GetSectionItems(ctx, "1")

// Fetch all episodes belonging to a TV show (flattened across all seasons)
episodes, _, err := client.Library.GetAllLeaves(ctx, "1001")
for _, ep := range episodes {
    fmt.Printf("S%02dE%02d: %s (Watched: %d times, TVDB: %s)\n",
        ep.ParentIndex, ep.Index, ep.Title, ep.ViewCount, ep.TVDB())
}
```

### 3. Active Playback Sessions (Live Monitoring)

```go
// Check who is watching what on the server in real time
sessions, _, err := client.Status.GetSessions(ctx)
if err != nil {
    log.Fatal(err)
}

for _, s := range sessions {
    user := "unknown"
    if s.User != nil {
        user = s.User.Title
    }
    state := "unknown"
    if s.Player != nil {
        state = s.Player.State
    }
    fmt.Printf("User %s is %s: %s (Progress: %d / %d ms)\n",
        user, state, s.Title, s.ViewOffset, s.Duration)
}
```

### 4. Real-Time Webhook Processing

```go
// In your HTTP server handler:
func plexWebhookHandler(w http.ResponseWriter, r *http.Request) {
    // ParseWebhook streams thumbnail attachments directly to io.Discard
    payload, err := plex.ParseWebhook(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Only process scrobble events for episodes
    if payload.IsScrobble() && payload.Metadata.Type == "episode" {
        fmt.Printf("User %s completed episode: %s S%02dE%02d (TVDB: %s)\n",
            payload.Account.Title,
            payload.Metadata.GrandparentTitle,
            payload.Metadata.ParentIndex,
            payload.Metadata.Index,
            payload.Metadata.TVDB(),
        )
    }

    w.WriteHeader(http.StatusOK)
}
```

### 5. Automatic Server Discovery (plex.tv)

```go
// Discover active server connection URI if local address is unknown or changed
activeURL := client.Discovery.EnsureActiveURL(ctx)
fmt.Printf("Active Plex Server URI: %s\n", activeURL)
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
