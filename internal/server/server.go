// Package server provides an HTTP API for the color science library, enabling
// the frontend to perform color conversions, delta-E calculations, and palette
// generation.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"colorsci/internal/colorconv"
	"colorsci/internal/colormetric"
	"colorsci/internal/colorparse"
)

// Server is the HTTP API server.
type Server struct {
	mux  *http.ServeMux
	addr string
}

// New creates a server.
func New(addr string) *Server {
	s := &Server{mux: http.NewServeMux(), addr: addr}
	s.mux.HandleFunc("/api/convert", s.handleConvert)
	s.mux.HandleFunc("/api/delta", s.handleDelta)
	s.mux.HandleFunc("/api/palette", s.handlePalette)
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.Handle("/", http.FileServer(http.Dir("frontend")))
	return s
}

// Handler returns the HTTP handler.
func (s *Server) Handler() http.Handler { return s.mux }

// ListenAndServe starts the server.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}

// ConvertRequest converts a color to multiple representations.
type ConvertRequest struct {
	Color string `json:"color"` // hex, rgb(), or named color
}

// ConvertResponse holds all color space representations.
type ConvertResponse struct {
	Hex string    `json:"hex"`
	RGB [3]uint8  `json:"rgb"`
	HSL [3]float64 `json:"hsl"`
	Lab [3]float64 `json:"lab"`
	LCh [3]float64 `json:"lch"`
	XYZ [3]float64 `json:"xyz"`
}

func (s *Server) handleConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	c, err := colorparse.Parse(req.Color)
	if err != nil {
		http.Error(w, "invalid color: "+err.Error(), http.StatusBadRequest)
		return
	}

	L, a, b := colorconv.ColorToLab(c)
	_, Ch, h := colorconv.LabToLCh(L, a, b)
	x, y, z := colorconv.ColorToXYZ(c)
	hue, sat, light := colorconv.RGBToHSL(c.R, c.G, c.B)

	resp := ConvertResponse{
		Hex: fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B),
		RGB: [3]uint8{c.R, c.G, c.B},
		HSL: [3]float64{hue, sat, light},
		Lab: [3]float64{L, a, b},
		LCh: [3]float64{L, Ch, h},
		XYZ: [3]float64{x, y, z},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeltaRequest computes color difference.
type DeltaRequest struct {
	Color1 string `json:"color1"`
	Color2 string `json:"color2"`
}

// DeltaResponse holds delta-E values.
type DeltaResponse struct {
	DeltaE76  float64 `json:"delta_e76"`
	DeltaE94  float64 `json:"delta_e94"`
	DeltaE00  float64 `json:"delta_e00"`
	Contrast  float64 `json:"contrast_ratio"`
}

func (s *Server) handleDelta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req DeltaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	c1, err := colorparse.Parse(req.Color1)
	if err != nil {
		http.Error(w, "invalid color1", http.StatusBadRequest)
		return
	}
	c2, err := colorparse.Parse(req.Color2)
	if err != nil {
		http.Error(w, "invalid color2", http.StatusBadRequest)
		return
	}

	lab1 := colormetric.LabOf(c1)
	lab2 := colormetric.LabOf(c2)

	resp := DeltaResponse{
		DeltaE76: colormetric.DeltaE76(lab1, lab2),
		DeltaE94: colormetric.DeltaE94(lab1, lab2),
		DeltaE00: colormetric.DeltaE2000(lab1, lab2),
		Contrast: colormetric.ContrastRatio(c1, c2),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PaletteRequest generates a color palette.
type PaletteRequest struct {
	BaseColor string `json:"base_color"`
	Count     int    `json:"count"`
	Mode      string `json:"mode"` // "analogous", "complement", "triadic"
}

func (s *Server) handlePalette(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req PaletteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	c, err := colorparse.Parse(req.BaseColor)
	if err != nil {
		http.Error(w, "invalid base_color", http.StatusBadRequest)
		return
	}
	if req.Count <= 0 {
		req.Count = 5
	}

	h, s2, l := colorconv.RGBToHSL(c.R, c.G, c.B)
	colors := make([]string, req.Count)
	for i := 0; i < req.Count; i++ {
		var hue float64
		switch req.Mode {
		case "complement":
			hue = h + 180*float64(i)/float64(req.Count)
		case "triadic":
			hue = h + 120*float64(i)/float64(req.Count-1+1)
		default: // analogous
			hue = h + 30*float64(i-req.Count/2)
		}
		for hue < 0 {
			hue += 360
		}
		for hue >= 360 {
			hue -= 360
		}
		pr, pg, pb := colorconv.HSLToRGB8(hue, s2, l)
		colors[i] = fmt.Sprintf("#%02x%02x%02x", pr, pg, pb)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"colors": colors})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}
