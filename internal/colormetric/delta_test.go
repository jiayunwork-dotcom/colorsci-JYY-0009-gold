package colormetric

import (
	"math"
	"testing"
)

func TestDeltaE76Identical(t *testing.T) {
	x := Lab{L: 50, A: 10, B: -10}
	if d := DeltaE76(x, x); d != 0 {
		t.Errorf("DeltaE76(x,x) = %v, want 0", d)
	}
}

func TestDeltaE76Known(t *testing.T) {
	x := Lab{L: 0, A: 0, B: 0}
	y := Lab{L: 100, A: 0, B: 0}
	if d := DeltaE76(x, y); math.Abs(d-100) > 1e-9 {
		t.Errorf("DeltaE76 black/white = %v, want 100", d)
	}
}

func TestDeltaE76Symmetric(t *testing.T) {
	x := Lab{L: 20, A: 30, B: 40}
	y := Lab{L: 80, A: -10, B: 5}
	if d1, d2 := DeltaE76(x, y), DeltaE76(y, x); d1 != d2 {
		t.Errorf("not symmetric: %v vs %v", d1, d2)
	}
}

func TestDeltaE94IdenticalZero(t *testing.T) {
	x := Lab{L: 50, A: 20, B: -30}
	if d := DeltaE94(x, x); d != 0 {
		t.Errorf("DeltaE94(x,x) = %v, want 0", d)
	}
}

func TestDeltaE94Asymmetric(t *testing.T) {
	// CIE94 uses the reference color's chroma, so swapping inputs can
	// change the result.
	x := Lab{L: 50, A: 0, B: 0}
	y := Lab{L: 50, A: 40, B: 0}
	d1 := DeltaE94(x, y)
	d2 := DeltaE94(y, x)
	if d1 == d2 {
		t.Errorf("expected asymmetry, got %v both ways", d1)
	}
	if d1 <= 0 || d2 <= 0 {
		t.Errorf("differences must be positive, got %v %v", d1, d2)
	}
}

func TestDeltaE2000Identical(t *testing.T) {
	x := Lab{L: 50, A: 2.5, B: 0}
	if d := DeltaE2000(x, x); d != 0 {
		t.Errorf("DeltaE2000(x,x) = %v, want 0", d)
	}
}

// Reference pairs from Sharma, Wu & Dalal, "The CIEDE2000
// Color-Difference Formula: Implementation Notes, Supplementary Test
// Data, and Mathematical Observations" (2005).
func TestDeltaE2000SharmaPairs(t *testing.T) {
	type pair struct {
		l1, a1, b1, l2, a2, b2, want float64
	}
	pairs := []pair{
		{50, 2.6772, -79.7751, 50, 0, -82.7485, 2.0425},
		{50, 3.1571, -77.2803, 50, 0, -82.7485, 2.8615},
		{50, 2.8361, -74.0200, 50, 0, -82.7485, 3.4412},
		{50, -1.3802, -84.2814, 50, 0, -82.7485, 1.0000},
		{50, -1.1848, -84.8006, 50, 0, -82.7485, 1.0000},
		{50, -0.9009, -85.5211, 50, 0, -82.7485, 1.0000},
		{50, 0, 0, 50, -1, 2, 2.3669},
		{50, -1, 2, 50, 0, 0, 2.3669},
		{50, 2.49, -0.001, 50, -2.49, 0.0009, 7.1792},
		{50, 2.5, 0, 73, 25, -18, 27.1492},
		{50, 2.5, 0, 61, -5, 29, 22.8977},
		{50, 2.5, 0, 56, -27, -3, 31.9030},
		{50, 2.5, 0, 58, 24, 15, 19.4535},
		{50, 2.5, 0, 50, 3.1736, 0.5854, 1.0000},
		{50, 2.5, 0, 50, 3.2972, 0, 1.0000},
		{50, 2.5, 0, 50, 1.8634, 0.5757, 1.0000},
		{50, 2.5, 0, 50, 3.2592, 0.3350, 1.0000},
		{60.2574, -34.0099, 36.2677, 60.4626, -34.1751, 39.4387, 1.2644},
		{63.0109, -31.0961, -5.8663, 62.8187, -29.7946, -4.0864, 1.2630},
		{22.7233, 20.0904, -46.6940, 23.0331, 14.9730, -42.5619, 2.0373},
		{6.7747, -0.2908, -2.4247, 5.8714, -0.0985, -2.2286, 0.6377},
		{2.0776, 0.0795, -1.1350, 0.9033, -0.0636, -0.5514, 0.9082},
	}
	for i, p := range pairs {
		x := Lab{L: p.l1, A: p.a1, B: p.b1}
		y := Lab{L: p.l2, A: p.a2, B: p.b2}
		got := DeltaE2000(x, y)
		if math.Abs(got-p.want) > 5e-3 {
			t.Errorf("pair %d: DeltaE2000 = %.6f, want %.4f", i+1, got, p.want)
		}
	}
}

func TestDeltaE2000Symmetric(t *testing.T) {
	x := Lab{L: 50, A: 2.5, B: 0}
	y := Lab{L: 73, A: 25, B: -18}
	if d1, d2 := DeltaE2000(x, y), DeltaE2000(y, x); math.Abs(d1-d2) > 1e-9 {
		t.Errorf("DeltaE2000 not symmetric: %v vs %v", d1, d2)
	}
}

func TestDeltaE2000Achromatic(t *testing.T) {
	// Colors below the chroma threshold must not produce NaN.
	x := Lab{L: 50, A: 0, B: 0}
	y := Lab{L: 51, A: 0, B: 0}
	d := DeltaE2000(x, y)
	if math.IsNaN(d) || d <= 0 {
		t.Errorf("achromatic DeltaE2000 = %v, want > 0", d)
	}
}
