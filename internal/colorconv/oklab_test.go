package colorconv

import (
	"colorsci/internal/colorparse"
	"math"
	"testing"
)

func TestOKLabRoundTrip(t *testing.T) {
	colors := []colorparse.Color{
		{R: 255, G: 0, B: 0},
		{R: 0, G: 255, B: 0},
		{R: 0, G: 0, B: 255},
		{R: 128, G: 128, B: 128},
	}
	for _, c := range colors {
		L, A, B := RGBToOKLab(c)
		back := OKLabToRGB(L, A, B)
		if absDiff(c.R, back.R) > 2 || absDiff(c.G, back.G) > 2 || absDiff(c.B, back.B) > 2 {
			t.Errorf("roundtrip failed: %v -> (%.3f,%.3f,%.3f) -> %v", c, L, A, B, back)
		}
	}
}

func TestOKLabLCh(t *testing.T) {
	L, A, B := 0.5, 0.1, -0.05
	l2, c2, h2 := OKLabToLCh(L, A, B)
	l3, a3, b3 := OKLChToLab(l2, c2, h2)
	if math.Abs(L-l3) > 1e-10 || math.Abs(A-a3) > 1e-10 || math.Abs(B-b3) > 1e-10 {
		t.Fatal("LCh roundtrip failed")
	}
}

func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}
