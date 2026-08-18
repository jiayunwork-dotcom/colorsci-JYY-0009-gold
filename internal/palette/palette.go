// Package palette implements color palette generation algorithms.
package palette

import (
	"math"

	"colorsci/internal/colorconv"
	"colorsci/internal/colorparse"
)

// Analogous generates n colors evenly spread around the base hue within angle degrees.
func Analogous(base colorparse.Color, n int, angle float64) []colorparse.Color {
	if n <= 0 {
		n = 5
	}
	h, s, l := colorconv.RGBToHSL(base.R, base.G, base.B)
	colors := make([]colorparse.Color, n)
	step := angle / float64(n-1)
	start := h - angle/2
	for i := 0; i < n; i++ {
		hue := start + step*float64(i)
		hue = normHue(hue)
		r, g, b := colorconv.HSLToRGB8(hue, s, l)
		colors[i] = colorparse.Color{R: r, G: g, B: b}
	}
	return colors
}

// Complementary returns the complementary color (180° opposite).
func Complementary(base colorparse.Color) colorparse.Color {
	h, s, l := colorconv.RGBToHSL(base.R, base.G, base.B)
	hue := normHue(h + 180)
	r, g, b := colorconv.HSLToRGB8(hue, s, l)
	return colorparse.Color{R: r, G: g, B: b}
}

// Triadic returns three colors evenly spaced 120° apart.
func Triadic(base colorparse.Color) [3]colorparse.Color {
	h, s, l := colorconv.RGBToHSL(base.R, base.G, base.B)
	var result [3]colorparse.Color
	for i := 0; i < 3; i++ {
		hue := normHue(h + float64(i)*120)
		r, g, b := colorconv.HSLToRGB8(hue, s, l)
		result[i] = colorparse.Color{R: r, G: g, B: b}
	}
	return result
}

// SplitComplementary returns the base plus two colors 150° on each side.
func SplitComplementary(base colorparse.Color) [3]colorparse.Color {
	h, s, l := colorconv.RGBToHSL(base.R, base.G, base.B)
	var result [3]colorparse.Color
	result[0] = base
	r1, g1, b1 := colorconv.HSLToRGB8(normHue(h+150), s, l)
	result[1] = colorparse.Color{R: r1, G: g1, B: b1}
	r2, g2, b2 := colorconv.HSLToRGB8(normHue(h+210), s, l)
	result[2] = colorparse.Color{R: r2, G: g2, B: b2}
	return result
}

// Monochromatic generates n shades of the same hue with varying lightness.
func Monochromatic(base colorparse.Color, n int) []colorparse.Color {
	if n <= 0 {
		n = 5
	}
	h, s, _ := colorconv.RGBToHSL(base.R, base.G, base.B)
	colors := make([]colorparse.Color, n)
	for i := 0; i < n; i++ {
		l := (float64(i) + 0.5) / float64(n)
		r, g, b := colorconv.HSLToRGB8(h, s, l)
		colors[i] = colorparse.Color{R: r, G: g, B: b}
	}
	return colors
}

// Warm generates n warm colors (red-orange-yellow range).
func Warm(n int) []colorparse.Color {
	colors := make([]colorparse.Color, n)
	for i := 0; i < n; i++ {
		hue := float64(i) * 60 / float64(n) // 0-60°
		r, g, b := colorconv.HSLToRGB8(hue, 0.8, 0.5)
		colors[i] = colorparse.Color{R: r, G: g, B: b}
	}
	return colors
}

// Cool generates n cool colors (blue-cyan-green range).
func Cool(n int) []colorparse.Color {
	colors := make([]colorparse.Color, n)
	for i := 0; i < n; i++ {
		hue := 180 + float64(i)*120/float64(n)
		r, g, b := colorconv.HSLToRGB8(hue, 0.7, 0.45)
		colors[i] = colorparse.Color{R: r, G: g, B: b}
	}
	return colors
}

func normHue(h float64) float64 {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	return h
}
