package colorparse

import "testing"

func TestFormatHexOpaque(t *testing.T) {
	c := Color{R: 0xff, G: 0x80, B: 0x00, A: 1}
	if got := FormatHex(c); got != "#ff8000" {
		t.Errorf("FormatHex = %q, want #ff8000", got)
	}
}

func TestFormatHexAlpha(t *testing.T) {
	c := Color{R: 0xff, G: 0x00, B: 0x00, A: 0.5}
	if got := FormatHex(c); got != "#ff000080" {
		t.Errorf("FormatHex(alpha) = %q, want #ff000080", got)
	}
}

func TestFormatRGBFunc(t *testing.T) {
	c := Color{R: 12, G: 34, B: 56, A: 1}
	if got := FormatRGB(c); got != "rgb(12 34 56)" {
		t.Errorf("FormatRGB = %q, want rgb(12 34 56)", got)
	}
	c.A = 0.25
	if got := FormatRGB(c); got != "rgb(12 34 56 / 0.25)" {
		t.Errorf("FormatRGB(alpha) = %q, want rgb(12 34 56 / 0.25)", got)
	}
}

func TestFormatRGBComma(t *testing.T) {
	c := Color{R: 1, G: 2, B: 3, A: 0.1}
	if got := FormatRGBComma(c); got != "rgba(1, 2, 3, 0.1)" {
		t.Errorf("FormatRGBComma = %q, want rgba(1, 2, 3, 0.1)", got)
	}
}

func TestFormatAlphaTrailingZeros(t *testing.T) {
	c := Color{A: 0.5000}
	if got := FormatRGB(c); got != "rgb(0 0 0 / 0.5)" {
		t.Errorf("trailing zero trimmed: %q", got)
	}
}

func TestStringAndOpaque(t *testing.T) {
	c := Color{R: 0x12, G: 0x34, B: 0x56, A: 1}
	if got := c.String(); got != "#123456" {
		t.Errorf("String = %q, want #123456", got)
	}
	if !c.Opaque() {
		t.Errorf("opaque color reported non-opaque")
	}
	if (Color{A: 0.99}).Opaque() {
		t.Errorf("alpha 0.99 reported opaque")
	}
}

func TestWithAlphaClamp(t *testing.T) {
	c := Color{R: 10, G: 20, B: 30, A: 1}
	low := c.WithAlpha(-3)
	if low.A != 0 {
		t.Errorf("WithAlpha(-3).A = %v, want 0", low.A)
	}
	high := c.WithAlpha(7)
	if high.A != 1 {
		t.Errorf("WithAlpha(7).A = %v, want 1", high.A)
	}
	if high.R != 10 {
		t.Errorf("WithAlpha mutated channels: R = %d", high.R)
	}
}

func TestParseFormatRoundTrip(t *testing.T) {
	inputs := []string{"#ff8000", "#f80", "#12345678", "rgb(1, 2, 3)", "hsl(210, 50%, 50%)"}
	for _, in := range inputs {
		c, err := Parse(in)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", in, err)
			continue
		}
		c2, err := Parse(FormatHex(c))
		if err != nil {
			t.Errorf("re-Parse(hex of %q) error: %v", in, err)
			continue
		}
		if c != c2 {
			t.Errorf("round trip %q: %v -> %v", in, c, c2)
		}
	}
}
