package colorparse

import (
	"fmt"
	"math"
	"strings"
)

// FormatHex renders the color in CSS hex syntax. Opaque colors use
// #rrggbb; colors with alpha < 1 use #rrggbbaa (alpha rounded to a
// byte).
func FormatHex(c Color) string {
	if c.A >= 1 {
		return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, roundByte(c.A))
}

func roundByte(v float64) byte {
	v = math.Max(0, math.Min(1, v))
	return byte(math.Round(v * 255))
}

// FormatRGB renders the color in modern rgb() syntax. Channels are the
// stored 8-bit values; alpha is included only when it is less than 1.
func FormatRGB(c Color) string {
	if c.A >= 1 {
		return fmt.Sprintf("rgb(%d %d %d)", c.R, c.G, c.B)
	}
	return fmt.Sprintf("rgb(%d %d %d / %s)", c.R, c.G, c.B, formatAlpha(c.A))
}

// FormatRGBComma renders the color in legacy rgba() syntax.
func FormatRGBComma(c Color) string {
	return fmt.Sprintf("rgba(%d, %d, %d, %s)", c.R, c.G, c.B, formatAlpha(c.A))
}

func formatAlpha(a float64) string {
	s := fmt.Sprintf("%.4f", a)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		s = "0"
	}
	return s
}

// MustParse parses s and panics on error. Use it only for inputs that
// are known to be valid (literals, validated configuration).
func MustParse(s string) Color {
	c, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return c
}

// String implements fmt.Stringer and returns FormatHex(c).
func (c Color) String() string {
	return FormatHex(c)
}

// WithAlpha returns a copy of the color with the given alpha. Values
// outside [0,1] are clamped.
func (c Color) WithAlpha(a float64) Color {
	c.A = math.Max(0, math.Min(1, a))
	return c
}

// Opaque reports whether the color is fully opaque (A >= 1).
func (c Color) Opaque() bool {
	return c.A >= 1
}
