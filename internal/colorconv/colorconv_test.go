package colorconv

import (
	"math"
	"testing"

	"colorsci/internal/colorparse"
)

func almost(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestLinearizeEncodeRoundTrip(t *testing.T) {
	// Boundary points are excluded: the sRGB thresholds 0.04045 and
	// 0.0031308 do not round-trip to bit-exact values.
	vals := []float64{0, 0.001, 0.01, 0.1, 0.5, 0.9, 1}
	for _, v := range vals {
		enc := Encode(Linearize(v))
		if !almost(enc, v, 1e-9) {
			t.Errorf("Encode(Linearize(%v)) = %v", v, enc)
		}
	}
}

func TestLinearizeBoundary(t *testing.T) {
	// At 0.04045 both branches coincide.
	if got := Linearize(0.04045); !almost(got, 0.04045/12.92, 1e-12) {
		t.Errorf("Linearize(0.04045) = %v", got)
	}
}

func TestEncodeClamps(t *testing.T) {
	if got := Encode(-1); got != 12.92*-1 {
		t.Errorf("Encode(-1) = %v", got)
	}
	if got := Linearize(2); got != 1 {
		t.Errorf("Linearize(2) = %v, want 1", got)
	}
}

func TestWhiteToXYZ(t *testing.T) {
	c := colorparse.Color{R: 255, G: 255, B: 255, A: 1}
	x, y, z := ColorToXYZ(c)
	if !almost(x, WhiteX, 1e-4) || !almost(y, WhiteY, 1e-4) || !almost(z, WhiteZ, 1e-4) {
		t.Errorf("white XYZ = (%v,%v,%v), want (%v,1,%v)", x, y, z, WhiteX, WhiteZ)
	}
}

func TestRedToXYZ(t *testing.T) {
	c := colorparse.Color{R: 255, G: 0, B: 0, A: 1}
	x, y, z := ColorToXYZ(c)
	if !almost(x, 0.4124564, 1e-4) || !almost(y, 0.2126729, 1e-4) || !almost(z, 0.0193339, 1e-4) {
		t.Errorf("red XYZ = (%v,%v,%v)", x, y, z)
	}
}

func TestXYZRGBRoundTrip(t *testing.T) {
	colors := []colorparse.Color{
		{R: 255, G: 255, B: 255, A: 1},
		{R: 0, G: 0, B: 0, A: 1},
		{R: 128, G: 64, B: 32, A: 1},
		{R: 10, G: 200, B: 90, A: 1},
	}
	for _, c := range colors {
		x, y, z := ColorToXYZ(c)
		back := XYZToColor(x, y, z)
		if back != c {
			t.Errorf("XYZ round trip %v -> %v", c, back)
		}
	}
}

func TestRedToLabKnown(t *testing.T) {
	c := colorparse.Color{R: 255, G: 0, B: 0, A: 1}
	l, a, b := ColorToLab(c)
	// Reference: sRGB (255,0,0) is Lab (53.233, 80.109, 67.220) under D65.
	if !almost(l, 53.233, 0.02) || !almost(a, 80.109, 0.02) || !almost(b, 67.220, 0.02) {
		t.Errorf("red Lab = (%v,%v,%v)", l, a, b)
	}
}

func TestWhiteToLab(t *testing.T) {
	c := colorparse.Color{R: 255, G: 255, B: 255, A: 1}
	l, a, b := ColorToLab(c)
	if !almost(l, 100, 1e-3) || !almost(a, 0, 1e-3) || !almost(b, 0, 1e-3) {
		t.Errorf("white Lab = (%v,%v,%v)", l, a, b)
	}
}

func TestBlackToLab(t *testing.T) {
	c := colorparse.Color{R: 0, G: 0, B: 0, A: 1}
	l, a, b := ColorToLab(c)
	if !almost(l, 0, 1e-3) || !almost(a, 0, 1e-3) || !almost(b, 0, 1e-3) {
		t.Errorf("black Lab = (%v,%v,%v)", l, a, b)
	}
}

func TestLabXYZRoundTrip(t *testing.T) {
	pairs := [][3]float64{
		{50, 0, 0}, {53.233, 80.109, 67.220}, {20, -30, 45}, {90, 5, -80}, {0, 0, 0}, {100, 0, 0},
	}
	for _, p := range pairs {
		x, y, z := LabToXYZ(p[0], p[1], p[2])
		l, a, b := XYZToLab(x, y, z)
		if !almost(l, p[0], 1e-6) || !almost(a, p[1], 1e-6) || !almost(b, p[2], 1e-6) {
			t.Errorf("Lab round trip %v -> (%v,%v,%v)", p, l, a, b)
		}
	}
}

func TestLChRoundTrip(t *testing.T) {
	labs := [][3]float64{{50, 30, 40}, {0, 0, 0}, {70, -20, -60}, {10, 100, 100}}
	for _, p := range labs {
		l, c, h := LabToLCh(p[0], p[1], p[2])
		l2, a2, b2 := LChToLab(l, c, h)
		if !almost(l2, p[0], 1e-9) || !almost(a2, p[1], 1e-9) || !almost(b2, p[2], 1e-9) {
			t.Errorf("LCh round trip %v -> (%v,%v,%v)", p, l2, a2, b2)
		}
	}
}

func TestLChAchromaticHueZero(t *testing.T) {
	_, c, h := LabToLCh(50, 0, 0)
	if c != 0 || h != 0 {
		t.Errorf("achromatic LCh = (%v,%v), want c=0 h=0", c, h)
	}
}

func TestLChHueRange(t *testing.T) {
	// Quadrant III should wrap to [0,360).
	_, _, h := LabToLCh(50, -10, -10)
	if h < 180 || h >= 360 {
		t.Errorf("hue for (-a,-b) = %v, want in [180,360)", h)
	}
}

func TestHSLRed(t *testing.T) {
	h, s, l := RGBToHSL(255, 0, 0)
	if !almost(h, 0, 1e-6) || !almost(s, 1, 1e-6) || !almost(l, 0.5, 1e-6) {
		t.Errorf("red HSL = (%v,%v,%v)", h, s, l)
	}
}

func TestHSLGray(t *testing.T) {
	_, s, l := RGBToHSL(128, 128, 128)
	if s != 0 {
		t.Errorf("gray saturation = %v, want 0", s)
	}
	if !almost(l, 128.0/255, 1e-6) {
		t.Errorf("gray lightness = %v", l)
	}
}

func TestHSLToRGB8(t *testing.T) {
	r, g, b := HSLToRGB8(120, 1, 0.5)
	if r != 0 || g != 255 || b != 0 {
		t.Errorf("hsl(120,1,0.5) -> %d,%d,%d, want 0,255,0", r, g, b)
	}
}

func TestHSLToRGB8NegativeHue(t *testing.T) {
	// -120 deg wraps to 240 deg, which is blue.
	r, g, b := HSLToRGB8(-120, 1, 0.5)
	if r != 0 || g != 0 || b != 255 {
		t.Errorf("hsl(-120,1,0.5) -> %d,%d,%d, want 0,0,255", r, g, b)
	}
}

func TestLabInGamut(t *testing.T) {
	inside := [][3]float64{{50, 0, 0}, {100, 0, 0}}
	// Red must be in gamut; use the exact conversion result rather than
	// a 3-decimal-rounded Lab triplet (rounding breaks the round trip).
	l, a, b := ColorToLab(colorparse.Color{R: 255, G: 0, B: 0, A: 1})
	inside = append(inside, [3]float64{l, a, b})
	for _, p := range inside {
		if !LabInGamut(p[0], p[1], p[2]) {
			t.Errorf("Lab%v should be in gamut", p)
		}
	}
	outside := [][3]float64{{50, 200, 0}, {0, 0, 100}, {50, 0, -200}}
	for _, p := range outside {
		if LabInGamut(p[0], p[1], p[2]) {
			t.Errorf("Lab%v should be out of gamut", p)
		}
	}
}

func TestClampLabToGamutInGamutUnchanged(t *testing.T) {
	l, a, b := ClampLabToGamut(50, 0, 0)
	if l != 50 || a != 0 || b != 0 {
		t.Errorf("in-gamut clamp changed color: (%v,%v,%v)", l, a, b)
	}
}

func TestClampLabToGamutOutOfGamut(t *testing.T) {
	original := [3]float64{50, 200, 0}
	l, a, b := ClampLabToGamut(original[0], original[1], original[2])
	if !LabInGamut(l, a, b) {
		t.Errorf("clamped color (%v,%v,%v) still out of gamut", l, a, b)
	}
	// Lightness preserved, hue preserved (chroma-only reduction).
	if !almost(l, original[0], 1e-6) {
		t.Errorf("lightness changed: %v -> %v", original[0], l)
	}
	_, c1, h1 := LabToLCh(original[0], original[1], original[2])
	_, c2, h2 := LabToLCh(l, a, b)
	if !almost(h1, h2, 1e-6) {
		t.Errorf("hue changed: %v -> %v", h1, h2)
	}
	if c2 >= c1 {
		t.Errorf("chroma not reduced: %v -> %v", c1, c2)
	}
}

func TestClampPreservesInGamut8(t *testing.T) {
	// After clamping, an 8-bit round trip should stay within tolerance.
	l, a, b := ClampLabToGamut(30, 120, -80)
	if !LabInGamut8(l, a, b, 2.0) {
		t.Errorf("clamped Lab (%v,%v,%v) not stable under 8-bit round trip", l, a, b)
	}
}
