package server

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	s := New(":0")
	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("got %d", w.Code)
	}
}

func TestConvert(t *testing.T) {
	s := New(":0")
	body := ConvertRequest{Color: "#ff0000"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/convert", bytes.NewReader(data))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("got %d: %s", w.Code, w.Body.String())
	}
	var resp ConvertResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.RGB[0] != 255 {
		t.Fatalf("expected R=255, got %d", resp.RGB[0])
	}
}

func TestDelta(t *testing.T) {
	s := New(":0")
	body := DeltaRequest{Color1: "#ff0000", Color2: "#00ff00"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/delta", bytes.NewReader(data))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("got %d", w.Code)
	}
	var resp DeltaResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.DeltaE76 <= 0 {
		t.Fatal("expected positive delta")
	}
}

func TestPalette(t *testing.T) {
	s := New(":0")
	body := PaletteRequest{BaseColor: "#3498db", Count: 5, Mode: "analogous"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/palette", bytes.NewReader(data))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("got %d", w.Code)
	}
}
