package vestaboard

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMessage(t *testing.T) {
	wantID := "abc-123"
	wantLayout := BoardLayout{
		{1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if got := r.Header.Get("X-Vestaboard-Token"); got != "test-token" {
			t.Errorf("expected token header 'test-token', got %q", got)
		}
		resp := map[string]any{
			"currentMessage": map[string]any{
				"id":     wantID,
				"layout": wantLayout,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := newWithURLs("test-token", srv.URL, srv.URL)
	result, err := client.GetMessage()
	if err != nil {
		t.Fatalf("GetMessage() error: %v", err)
	}
	if result.ID != wantID {
		t.Errorf("ID: got %q, want %q", result.ID, wantID)
	}
	if len(result.Layout) != len(wantLayout) {
		t.Fatalf("layout rows: got %d, want %d", len(result.Layout), len(wantLayout))
	}
	if result.Layout[0][0] != wantLayout[0][0] || result.Layout[0][2] != wantLayout[0][2] {
		t.Errorf("layout mismatch: got %v, want %v", result.Layout[0], wantLayout[0])
	}
}

func TestSendText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if got := r.Header.Get("X-Vestaboard-Token"); got != "test-token" {
			t.Errorf("expected token header 'test-token', got %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		json.Unmarshal(body, &payload)
		if payload["text"] != "hello" {
			t.Errorf("expected text 'hello', got %v", payload["text"])
		}
		if _, ok := payload["forced"]; ok {
			t.Errorf("forced should not be set")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "id": "msg-1", "created": 1681154452865})
	}))
	defer srv.Close()

	client := newWithURLs("test-token", srv.URL, srv.URL)
	result, err := client.SendText("hello", false)
	if err != nil {
		t.Fatalf("SendText() error: %v", err)
	}
	if result.Status != "ok" {
		t.Errorf("status: got %q, want 'ok'", result.Status)
	}
	if result.ID != "msg-1" {
		t.Errorf("id: got %q, want 'msg-1'", result.ID)
	}
}

func TestSendTextForced(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		json.Unmarshal(body, &payload)
		if payload["forced"] != true {
			t.Errorf("expected forced=true, got %v", payload["forced"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "id": "msg-2", "created": 1681154452865})
	}))
	defer srv.Close()

	client := newWithURLs("test-token", srv.URL, srv.URL)
	result, err := client.SendText("urgent", true)
	if err != nil {
		t.Fatalf("SendText(forced) error: %v", err)
	}
	if result.ID != "msg-2" {
		t.Errorf("id: got %q, want 'msg-2'", result.ID)
	}
}
