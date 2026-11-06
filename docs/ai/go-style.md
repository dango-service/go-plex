# Go Style Rules  -  go-plex

Coding conventions for `go-plex`.

---

## 1. Error Handling

- Always wrap errors using `%w`:
  ```go
  return nil, nil, fmt.Errorf("plex: %s: %w", op, err)
  ```
- Expose typed errors and predicates:
  ```go
  var (
      ErrUnauthorized = errors.New("plex: unauthorized (401)")
      ErrNotFound     = errors.New("plex: not found (404)")
      ErrBadRequest   = errors.New("plex: bad request (400)")
      ErrIgnoredEvent = errors.New("plex: webhook event ignored")
  )
  ```

---

## 2. Context Cancellation

- `context.Context` must be the first parameter in all service methods performing network I/O.
- Check context cancellation early:
  ```go
  if err := ctx.Err(); err != nil {
      return nil, nil, fmt.Errorf("plex: %s: %w", op, err)
  }
  ```

---

## 3. Streaming Bounded Deserialization

- Always use `io.LimitReader`:
  ```go
  limitedBody := io.LimitReader(resp.Body, maxBodyBytes)
  if err := json.NewDecoder(limitedBody).Decode(v); err != nil {
      return nil, fmt.Errorf("plex: decode response: %w", err)
  }
  ```

---

## 4. Pure CGO_ENABLED=0

- No cgo dependencies.
- Library must compile on Windows, Linux, and macOS without a C compiler.
