package plex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, v any) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("plex: context error: %w", err)
	}

	urlStr := path
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = c.baseURL + "/" + strings.TrimLeft(path, "/")
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	if err != nil {
		return nil, fmt.Errorf("plex: create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("X-Plex-Token", c.token)
	}
	if c.clientID != "" {
		req.Header.Set("X-Plex-Client-Identifier", c.clientID)
	}
	if c.product != "" {
		req.Header.Set("X-Plex-Product", c.product)
	}
	if c.version != "" {
		req.Header.Set("X-Plex-Version", c.version)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("plex: request canceled: %w", ctxErr)
		}
		return nil, fmt.Errorf("plex: execute request: %w", err)
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return response, c.handleErrorResponse(resp)
	}

	if v != nil {
		limitedReader := io.LimitReader(resp.Body, c.maxBodyBytes)
		if err := json.NewDecoder(limitedReader).Decode(v); err != nil && err != io.EOF {
			return response, fmt.Errorf("plex: decode response: %w", err)
		}
	}

	return response, nil
}

func (c *Client) handleErrorResponse(resp *http.Response) error {
	limitedReader := io.LimitReader(resp.Body, c.maxBodyBytes)
	errResp := &ErrorResponse{
		StatusCode: resp.StatusCode,
	}

	_ = json.NewDecoder(limitedReader).Decode(errResp)

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		errResp.Err = ErrUnauthorized
	case http.StatusNotFound:
		errResp.Err = ErrNotFound
	case http.StatusBadRequest:
		errResp.Err = ErrBadRequest
	default:
		errResp.Err = fmt.Errorf("plex: HTTP %d", resp.StatusCode)
	}

	return errResp
}
