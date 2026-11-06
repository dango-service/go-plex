# Testing Strategy  -  go-plex

Deterministic testing strategy for `go-plex`.

---

## 1. Zero Network Policy

- All unit tests must communicate exclusively with `net/http/httptest.Server`.
- Real network calls to local PMS servers or `plex.tv` are strictly prohibited during tests.

---

## 2. Table-Driven Tests & Fixtures

- Tests must use `t.Run` for table-driven scenarios covering:
  - PMS MediaContainer responses with realistic JSON fixtures in `testdata/`.
  - Webhook multipart requests with both thumbnail images and JSON payloads.
  - Context cancellation, timeout handling, and malformed responses.
  - Flexible ID decoding (both JSON numbers and strings).

---

## 3. Test Coverage Goals

- Core transport, flexible ID unmarshaling, and webhook thumbnail discarding: 100% coverage.
- All library and status endpoints: covered by positive and negative tests.
