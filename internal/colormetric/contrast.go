package colormetric

import (
	"math"

	"colorsci/internal/colorconv"
	"colorsci/internal/colorparse"
)

// RelativeLuminance computes the WCAG 2.x relative luminance of an
// sRGB8 color: 0 for black, 1 for white.
func RelativeLuminance(c colorparse.Color) float64 {
	r, g, b := colorconv.RGBToLinear(c.R, c.G, c.B)
	return 0.2126*r + 0.7152*g + 0.0722*b
}

// ContrastRatio computes the WCAG contrast ratio between two colors in
// the range [1,21]. The result is symmetric.
func ContrastRatio(c1, c2 colorparse.Color) float64 {
	l1 := RelativeLuminance(c1)
	l2 := RelativeLuminance(c2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// Levels reports which WCAG 2.1 success criteria a contrast ratio
// satisfies. large refers to "large text" (>=18pt or >=14pt bold).
type Levels struct {
	AANormal  bool // 4.5:1
	AALarge   bool // 3.0:1
	AAANormal bool // 7.0:1
	AAALarge  bool // 4.5:1
}

// WCAGLevels classifies a contrast ratio against the four WCAG 2.1
// thresholds.
func WCAGLevels(ratio float64) Levels {
	return Levels{
		AANormal:  ratio >= 4.5,
		AALarge:   ratio >= 3.0,
		AAANormal: ratio >= 7.0,
		AAALarge:  ratio >= 4.5,
	}
}

// Passes reports whether the pair meets the given criteria.
func (l Levels) Passes(large, aaa bool) bool {
	if large {
		if aaa {
			return l.AAALarge
		}
		return l.AALarge
	}
	if aaa {
		return l.AAANormal
	}
	return l.AANormal
}

// BestPair finds the pair of colors from a palette (first element is
// the foreground candidate set, second the background set) with the
// highest contrast ratio. It returns the two colors and the ratio. When
// both sets are non-empty the result is deterministic.
func BestPair(foreground, background []colorparse.Color) (colorparse.Color, colorparse.Color, float64) {
	best := -1.0
	var bf, bb colorparse.Color
	for _, f := range foreground {
		for _, b := range background {
			r := ContrastRatio(f, b)
			if r > best {
				best = r
				bf, bb = f, b
			}
		}
	}
	if best < 0 {
		return colorparse.Color{}, colorparse.Color{}, math.Inf(-1)
	}
	return bf, bb, best
}
