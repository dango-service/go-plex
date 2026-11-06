# Plex Media Server and Webhook Technical Reference

Technical specification for Plex Media Server (PMS) HTTP API and Plex Webhooks.

---

## 1. Authentication & Headers

- PMS authentication: `X-Plex-Token: <token>`.
- Client identification: `X-Plex-Client-Identifier: <client_uuid>`.
- Output format: `Accept: application/json`.
- Product header: `X-Plex-Product: go-plex`.
- Version header: `X-Plex-Version: 0.1.0`.

---

## 2. Media Container Structure

Plex wraps all JSON array responses in a top-level `MediaContainer` object:
```json
{
  "MediaContainer": {
    "size": 1,
    "Metadata": [ ... ],
    "Directory": [ ... ]
  }
}
```

---

## 3. Key Endpoints Overview

- `GET /library/sections`: List all libraries (Movies, TV Shows, Anime).
- `GET /library/sections/{sectionKey}/all`: List all items in a library section.
- `GET /library/metadata/{ratingKey}`: Detailed item metadata.
- `GET /library/metadata/{ratingKey}/allLeaves`: All child episodes belonging to a TV show or season.
- `GET /status/sessions`: Currently active streaming sessions across all users.
- `GET /status/sessions/history/all`: Server watch history.
- `GET https://plex.tv/api/v2/resources?includeHttps=1&includeRelay=1`: Discover active server connection URIs.

---

## 4. Webhook Payload Format

- Encoded as `multipart/form-data` with two parts:
  1. `thumb`: JPEG/PNG binary image (discarded on the fly).
  2. `payload`: JSON string representing the watch event (`media.scrobble`, `media.play`, `media.pause`, `media.stop`, `media.rate`).
