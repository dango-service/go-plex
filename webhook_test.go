package plex

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseWebhook_Success(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// 1. Add binary thumbnail part
	thumbPart, err := writer.CreateFormFile("thumb", "thumb.jpg")
	if err != nil {
		t.Fatalf("failed to create thumb part: %v", err)
	}
	_, _ = thumbPart.Write([]byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46})

	// 2. Add JSON payload part
	payloadJSON := `{
		"event": "media.scrobble",
		"user": true,
		"owner": true,
		"Account": {
			"id": 1,
			"title": "hicham"
		},
		"Server": {
			"title": "Tower-PMS",
			"uuid": "server-uuid-1234"
		},
		"Player": {
			"local": true,
			"title": "Plex Web",
			"uuid": "player-uuid"
		},
		"Metadata": {
			"ratingKey": "5001",
			"type": "episode",
			"title": "The Journey's End",
			"index": 1,
			"parentIndex": 1,
			"grandparentTitle": "Sousou no Frieren",
			"Guid": [
				{"id": "tvdb://404080"},
				{"id": "tmdb://209867"}
			]
		}
	}`
	if err := writer.WriteField("payload", payloadJSON); err != nil {
		t.Fatalf("failed to write payload field: %v", err)
	}
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/webhook", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	payload, err := ParseWebhook(req)
	if err != nil {
		t.Fatalf("unexpected error parsing webhook: %v", err)
	}

	if !payload.IsScrobble() {
		t.Errorf("expected IsScrobble() to be true")
	}
	if payload.IsPlay() || payload.IsStop() || payload.IsPause() || payload.IsRate() {
		t.Errorf("expected other event predicates to be false")
	}
	if payload.Account.Title != "hicham" || payload.Server.Title != "Tower-PMS" {
		t.Errorf("unexpected account or server: %+v, %+v", payload.Account, payload.Server)
	}
	if payload.Metadata.Title != "The Journey's End" || payload.Metadata.TVDB() != "404080" {
		t.Errorf("unexpected metadata: %+v", payload.Metadata)
	}
	if payload.Metadata.Index != 1 || payload.Metadata.ParentIndex != 1 {
		t.Errorf("unexpected indices: %+v", payload.Metadata)
	}
}

func TestParseWebhook_Validation(t *testing.T) {
	// 1. Nil request
	if _, err := ParseWebhook(nil); err == nil {
		t.Fatal("expected error on nil request, got nil")
	}

	// 2. Invalid method
	getReq := httptest.NewRequest(http.MethodGet, "/webhook", nil)
	if _, err := ParseWebhook(getReq); err == nil {
		t.Fatal("expected error on GET request, got nil")
	}

	// 3. Not multipart
	plainReq := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString("hello"))
	if _, err := ParseWebhook(plainReq); err == nil {
		t.Fatal("expected error on non-multipart request, got nil")
	}

	// 4. Missing payload part
	var emptyBody bytes.Buffer
	w := multipart.NewWriter(&emptyBody)
	_ = w.WriteField("unrelated", "value")
	_ = w.Close()

	missingReq := httptest.NewRequest(http.MethodPost, "/webhook", &emptyBody)
	missingReq.Header.Set("Content-Type", w.FormDataContentType())
	if _, err := ParseWebhook(missingReq); err == nil {
		t.Fatal("expected error on missing payload part, got nil")
	}
}

func TestWebhookPayload_Predicates(t *testing.T) {
	scrobble := WebhookPayload{Event: EventMediaScrobble}
	play := WebhookPayload{Event: EventMediaPlay}
	pause := WebhookPayload{Event: EventMediaPause}
	stop := WebhookPayload{Event: EventMediaStop}
	rate := WebhookPayload{Event: EventMediaRate}

	if !scrobble.IsScrobble() {
		t.Error("expected IsScrobble() true")
	}
	if !play.IsPlay() {
		t.Error("expected IsPlay() true")
	}
	if !pause.IsPause() {
		t.Error("expected IsPause() true")
	}
	if !stop.IsStop() {
		t.Error("expected IsStop() true")
	}
	if !rate.IsRate() {
		t.Error("expected IsRate() true")
	}
}
