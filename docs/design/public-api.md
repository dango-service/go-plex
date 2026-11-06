# Public API Contract  -  go-plex

Public interface contract and signatures for `go-plex`.

---

## 1. Client Constructor & Options

```go
type Client struct {
    Library   *LibraryService
    Status    *StatusService
    Discovery *DiscoveryService
}

func NewClient(baseURL, token string, opts ...ClientOption) *Client

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption
func WithClientIdentifier(clientID string) ClientOption
func WithMaxBodyBytes(maxBytes int64) ClientOption
```

---

## 2. Service Signatures

### LibraryService
```go
func (s *LibraryService) GetSections(ctx context.Context) ([]Section, *Response, error)
func (s *LibraryService) GetSectionItems(ctx context.Context, sectionKey string) ([]MetadataItem, *Response, error)
func (s *LibraryService) GetMetadata(ctx context.Context, ratingKey string) (*MetadataItem, *Response, error)
func (s *LibraryService) GetAllLeaves(ctx context.Context, ratingKey string) ([]MetadataItem, *Response, error)
```

### StatusService
```go
func (s *StatusService) GetSessions(ctx context.Context) ([]SessionItem, *Response, error)
func (s *StatusService) GetHistory(ctx context.Context, minViewedAt int64) ([]HistoryItem, *Response, error)
```

### DiscoveryService
```go
func (s *DiscoveryService) DiscoverServers(ctx context.Context, token string) ([]ServerDevice, error)
```

### Webhook Parser
```go
func ParseWebhook(r *http.Request) (*WebhookPayload, error)
```
