package colorparse

import (
	"errors"
	"testing"
)

func TestParseHexLonghand(t *testing.T) {
	c, err := Parse("#ff8040")
	if err != nil {
		t.Fatalf("Parse(#ff8040) error: %v", err)
	}
	if c.R != 0xff || c.G != 0x80 || c.B != 0x40 {
		t.Errorf("Parse(#ff8040) = %02x%02x%02x, want ff8040", c.R, c.G, c.B)
	}
	if c.A != 1 {
		t.Errorf("alpha = %v, want 1", c.A)
	}
}

func TestParseHexShorthand(t *testing.T) {
	c, err := Parse("#f80")
	if err != nil {
		t.Fatalf("Parse(#f80) error: %v", err)
	}
	if c.R != 0xff || c.G != 0x88 || c.B != 0x00 {
		t.Errorf("Parse(#f80) = %02x%02x%02x, want ff8800", c.R, c.G, c.B)
	}
}

func TestParseHexAlpha(t *testing.T) {
	c, err := Parse("#ff000080")
	if err != nil {
		t.Fatalf("Parse(#ff000080) error: %v", err)
	}
	if c.R != 0xff || c.A != 0x80/255.0 {
		t.Errorf("got %v alpha %v, want ff0000 alpha %v", c, c.A, 0x80/255.0)
	}
	c4, err := Parse("#f008")
	if err != nil {
		t.Fatalf("Parse(#f008) error: %v", err)
	}
	if c4.A != 0x88/255.0 {
		t.Errorf("shorthand alpha = %v, want %v", c4.A, 0x88/255.0)
	}
}

func TestParseHexErrors(t *testing.T) {
	bad := []string{"#", "#12", "#12345", "#1234567", "#gg0000", "#12 345"}
	for _, s := range bad {
		if _, err := Parse(s); !errors.Is(err, ErrBadSyntax) {
			t.Errorf("Parse(%q) err = %v, want ErrBadSyntax", s, err)
		}
	}
}

func TestParseRGBFunc(t *testing.T) {
	c, err := Parse("rgb(255, 0, 128)")
	if err != nil {
		t.Fatalf("Parse(rgb) error: %v", err)
	}
	if c.R != 255 || c.G != 0 || c.B != 128 {
		t.Errorf("rgb(255,0,128) = %v", c)
	}
}

func TestParseRGBFuncPercent(t *testing.T) {
	c, err := Parse("rgb(100%, 0%, 50%)")
	if err != nil {
		t.Fatalf("Parse(percent) error: %v", err)
	}
	if c.R != 255 || c.G != 0 || c.B != 128 {
		t.Errorf("rgb(100%%,0%%,50%%) = %02x%02x%02x, want ff0080", c.R, c.G, c.B)
	}
}

func TestParseRGBFuncMixedUnitsRejected(t *testing.T) {
	if _, err := Parse("rgb(255, 0%, 50%)"); !errors.Is(err, ErrBadSyntax) {
		t.Errorf("mixed units err = %v, want ErrBadSyntax", err)
	}
}

func TestParseRGBFuncModernSlash(t *testing.T) {
	c, err := Parse("rgb(255 0 0 / 50%)")
	if err != nil {
		t.Fatalf("Parse(modern) error: %v", err)
	}
	if c.A != 0.5 {
		t.Errorf("alpha = %v, want 0.5", c.A)
	}
}

func TestParseRGBFuncOutOfRange(t *testing.T) {
	if _, err := Parse("rgb(300, 0, 0)"); err != nil {
		t.Fatalf("300 should clamp, got error %v", err)
	}
	c, _ := Parse("rgb(300, 0, 0)")
	if c.R != 255 {
		t.Errorf("clamped R = %d, want 255", c.R)
	}
}

func TestParseAlphaOutOfRange(t *testing.T) {
	if _, err := Parse("rgba(0, 0, 0, 1.5)"); !errors.Is(err, ErrBadAlpha) {
		t.Errorf("err = %v, want ErrBadAlpha", err)
	}
}

func TestParseHSLFunc(t *testing.T) {
	c, err := Parse("hsl(120, 100%, 50%)")
	if err != nil {
		t.Fatalf("Parse(hsl) error: %v", err)
	}
	if c.R != 0 || c.G != 255 || c.B != 0 {
		t.Errorf("hsl(120,100%%,50%%) = %02x%02x%02x, want 00ff00", c.R, c.G, c.B)
	}
}

func TestParseHSLFuncUnits(t *testing.T) {
	c, err := Parse("hsl(2.0944rad, 100%, 50%)")
	if err != nil {
		t.Fatalf("Parse(rad) error: %v", err)
	}
	if c.R != 0 || c.G != 255 || c.B != 0 {
		t.Errorf("hsl(2.0944rad) = %02x%02x%02x, want 00ff00 (120deg)", c.R, c.G, c.B)
	}
}

func TestParseHSLFuncHueWrap(t *testing.T) {
	c1, err := Parse("hsl(720, 100%, 50%)")
	if err != nil {
		t.Fatalf("Parse(720) error: %v", err)
	}
	c2, _ := Parse("hsl(360, 100%, 50%)")
	if c1 != c2 {
		t.Errorf("hsl(720) = %v, hsl(360) = %v, want equal", c1, c2)
	}
}

func TestParseHSLFuncMissingPercent(t *testing.T) {
	if _, err := Parse("hsl(120, 1, 0.5)"); !errors.Is(err, ErrBadSyntax) {
		t.Errorf("err = %v, want ErrBadSyntax (non-percent s/l)", err)
	}
}

func TestParseEmpty(t *testing.T) {
	if _, err := Parse("   "); !errors.Is(err, ErrEmpty) {
		t.Errorf("err = %v, want ErrEmpty", err)
	}
}

func TestParseCaseInsensitive(t *testing.T) {
	c, err := Parse("#FF0000")
	if err != nil {
		t.Fatalf("Parse(upper) error: %v", err)
	}
	if c.R != 255 {
		t.Errorf("upper hex R = %d, want 255", c.R)
	}
	n, err := Parse("RED")
	if err != nil {
		t.Fatalf("Parse(RED) error: %v", err)
	}
	if n.R != 255 || n.G != 0 {
		t.Errorf("RED = %v, want ff0000", n)
	}
}
