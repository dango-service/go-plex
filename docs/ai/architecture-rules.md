# Architecture Rules  -  go-plex

Architectural invariants for `go-plex`. Every contribution must uphold these constraints.

---

## 1. Zero External Dependencies

- `go-plex` must depend ONLY on the Go standard library.
- No third-party Plex wrappers, multipart parsers, or HTTP clients.
- Standard packages only: `net/http`, `mime/multipart`, `encoding/json`, `context`, `time`, `net/url`, `fmt`, `errors`, `io`.

---

## 2. Stateless Client & Authentication

- The root `Client` holds references to modular domain services (`Library`, `Status`, `Discovery`).
- Authentication is passed via HTTP header `X-Plex-Token: <token>`.
- Client identification header: `X-Plex-Client-Identifier: <id>` (defaulting to a stable client ID unless overridden).
- All requests must include `Accept: application/json`.

---

## 3. Streaming Deserialization & Bounded Memory

- Deserialization MUST use `json.NewDecoder` on `io.LimitReader(resp.Body, maxBodyBytes)`.
- Default body ceiling: 10 MiB (`10 * 1024 * 1024` bytes).
- Webhook multipart processing must discard thumbnail parts immediately using `io.Copy(io.Discard, part)`.

---

## 4. Flexible ID Resilience

- Plex Media Server represents IDs as either numbers or strings across different endpoints and versions.
- Types like `RatingKey` and `LibrarySectionID` must unmarshal gracefully whether the JSON value is an integer or a string.

---

## 5. Security & Credential Protection

- `Client` credentials and tokens must never be logged in raw form.
- Pass tokens via HTTP headers rather than URL query parameters to avoid proxy log leakage.
