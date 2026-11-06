# Scope and Bounded Context  -  go-plex

Scope of supported features and boundaries for `go-plex`.

---

## 1. Supported Features

- Querying library sections, library items, item metadata, and episode leaves (`allLeaves`).
- External GUID resolution helpers (`TVDB()`, `TMDB()`, `IMDb()`, `AniDB()`).
- Parsing genres, collections, and tags from media items.
- Active streaming playback session inspection (`/status/sessions`).
- Playback history inspection (`/status/sessions/history/all`).
- Memory-safe multipart streaming webhook parsing.
- Server connection discovery via `plex.tv`.

---

## 2. Excluded Features

- Local SQLite tracking database operations (responsibility of Dango).
- Tracker mutations (AniList, Kitsu, Trakt, Simkl, MyAnimeList).
- Media transcoding control or media file downloading.
