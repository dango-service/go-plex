package plex

import "net/http"

// Response wraps http.Response for Plex requests.
type Response struct {
	*http.Response
}
