package vestaboard

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFormatMessage(t *testing.T) {
	wantLayout := BoardLayout{
		{CharH, CharE, CharL, CharL, CharO, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/format" {
			t.Errorf("expected path /format, got %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		json.Unmarshal(body, &payload)
		if payload["message"] != "HELLO" {
			t.Errorf("expected message 'HELLO', got %v", payload["message"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wantLayout)
	}))
	defer srv.Close()

	client := newWithURLs("test-token", srv.URL, srv.URL)
	layout, err := client.FormatMessage("HELLO")
	if err != nil {
		t.Fatalf("FormatMessage() error: %v", err)
	}
	if len(layout) != len(wantLayout) {
		t.Fatalf("layout rows: got %d, want %d", len(layout), len(wantLayout))
	}
	if layout[0][0] != CharH || layout[0][4] != CharO {
		t.Errorf("layout[0]: got %v, want %v", layout[0], wantLayout[0])
	}
}

func TestGetTransition(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/transition" {
			t.Errorf("expected path /transition, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"transition": "wave", "transitionSpeed": "gentle"})
	}))
	defer srv.Close()

	client := newWithURLs("test-token", srv.URL, srv.URL)
	result, err := client.GetTransition()
	if err != nil {
		t.Fatalf("GetTransition() error: %v", err)
	}
	if result.Transition != TransitionWave {
		t.Errorf("transition: got %q, want %q", result.Transition, TransitionWave)
	}
	if result.TransitionSpeed != SpeedGentle {
		t.Errorf("speed: got %q, want %q", result.TransitionSpeed, SpeedGentle)
	}
}

func TestSetTransition(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		json.Unmarshal(body, &payload)
		if payload["transition"] != "curtain" {
			t.Errorf("transition: got %v, want 'curtain'", payload["transition"])
		}
		if payload["transitionSpeed"] != "fast" {
			t.Errorf("transitionSpeed: got %v, want 'fast'", payload["transitionSpeed"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"transition": "curtain", "transitionSpeed": "fast"})
	}))
	defer srv.Close()

	client := newWithURLs("test-token", srv.URL, srv.URL)
	result, err := client.SetTransition(TransitionCurtain, SpeedFast)
	if err != nil {
		t.Fatalf("SetTransition() error: %v", err)
	}
	if result.Transition != TransitionCurtain {
		t.Errorf("transition: got %q, want %q", result.Transition, TransitionCurtain)
	}
	if result.TransitionSpeed != SpeedFast {
		t.Errorf("speed: got %q, want %q", result.TransitionSpeed, SpeedFast)
	}
}

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

func TestSendCharacters(t *testing.T) {
	wantLayout := BoardLayout{
		{CharH, CharI, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		json.Unmarshal(body, &payload)
		if _, ok := payload["characters"]; !ok {
			t.Errorf("expected 'characters' key in payload")
		}
		if _, ok := payload["text"]; ok {
			t.Errorf("'text' key should not be present for SendCharacters")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "id": "msg-3", "created": 1681154452865})
	}))
	defer srv.Close()

	client := newWithURLs("test-token", srv.URL, srv.URL)
	result, err := client.SendCharacters(wantLayout, false)
	if err != nil {
		t.Fatalf("SendCharacters() error: %v", err)
	}
	if result.Status != "ok" {
		t.Errorf("status: got %q, want 'ok'", result.Status)
	}
	if result.ID != "msg-3" {
		t.Errorf("id: got %q, want 'msg-3'", result.ID)
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
