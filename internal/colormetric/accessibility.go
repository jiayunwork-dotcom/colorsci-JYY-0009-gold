package colormetric

import "colorsci/internal/colorparse"

// WCAGLevel represents WCAG 2.1 conformance levels.
type WCAGLevel string

const (
	WCAGFail WCAGLevel = "Fail"
	WCAGAA   WCAGLevel = "AA"
	WCAGAAA  WCAGLevel = "AAA"
)

// WCAGCheck checks a foreground/background pair against WCAG 2.1 criteria.
func WCAGCheck(fg, bg colorparse.Color, isLargeText bool) WCAGLevel {
	ratio := ContrastRatio(fg, bg)
	if isLargeText {
		if ratio >= 4.5 {
			return WCAGAAA
		}
		if ratio >= 3.0 {
			return WCAGAA
		}
	} else {
		if ratio >= 7.0 {
			return WCAGAAA
		}
		if ratio >= 4.5 {
			return WCAGAA
		}
	}
	return WCAGFail
}

// SuggestForeground finds a suitable foreground color (black or white)
// for the given background.
func SuggestForeground(bg colorparse.Color) colorparse.Color {
	white := colorparse.Color{R: 255, G: 255, B: 255}
	black := colorparse.Color{R: 0, G: 0, B: 0}
	if ContrastRatio(white, bg) > ContrastRatio(black, bg) {
		return white
	}
	return black
}

// IsReadable reports whether fg on bg meets minimum contrast (4.5:1).
func IsReadable(fg, bg colorparse.Color) bool {
	return ContrastRatio(fg, bg) >= 4.5
}

// ContrastDescription returns a human-readable description of the contrast.
func ContrastDescription(ratio float64) string {
	switch {
	case ratio >= 7:
		return "excellent"
	case ratio >= 4.5:
		return "good"
	case ratio >= 3:
		return "fair (large text only)"
	default:
		return "poor"
	}
}
