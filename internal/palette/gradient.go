package palette

import (
	"colorsci/internal/colorconv"
	"colorsci/internal/colorparse"
)

// Gradient generates a linear gradient of n colors between start and end
// by interpolating in HSL space. Both endpoints are included in the result.
// If n < 2 it defaults to 2.
func Gradient(start, end colorparse.Color, n int) []colorparse.Color {
	if n < 2 {
		n = 2
	}
	h1, s1, l1 := colorconv.RGBToHSL(start.R, start.G, start.B)
	h2, s2, l2 := colorconv.RGBToHSL(end.R, end.G, end.B)

	// 取最短色相路径
	dh := h2 - h1
	if dh > 180 {
		dh -= 360
	} else if dh < -180 {
		dh += 360
	}

	colors := make([]colorparse.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		h := normHue(h1 + dh*t)
		s := s1 + (s2-s1)*t
		l := l1 + (l2-l1)*t
		r, g, b := colorconv.HSLToRGB8(h, s, l)
		colors[i] = colorparse.Color{R: r, G: g, B: b, A: 1}
	}
	return colors
}

// GradientRGB generates a linear gradient by interpolating in RGB space.
func GradientRGB(start, end colorparse.Color, n int) []colorparse.Color {
	if n < 2 {
		n = 2
	}
	colors := make([]colorparse.Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		r := float64(start.R) + (float64(end.R)-float64(start.R))*t
		g := float64(start.G) + (float64(end.G)-float64(start.G))*t
		b := float64(start.B) + (float64(end.B)-float64(start.B))*t
		colors[i] = colorparse.Color{
			R: uint8(r + 0.5),
			G: uint8(g + 0.5),
			B: uint8(b + 0.5),
			A: 1,
		}
	}
	return colors
}
