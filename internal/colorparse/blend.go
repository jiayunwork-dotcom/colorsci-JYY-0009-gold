package colorparse

import "math"

// Blend linearly interpolates between two colors by factor t in [0,1].
func Blend(a, b Color, t float64) Color {
	if t <= 0 {
		return a
	}
	if t >= 1 {
		return b
	}
	return Color{
		R: lerp8(a.R, b.R, t),
		G: lerp8(a.G, b.G, t),
		B: lerp8(a.B, b.B, t),
	}
}

// BlendN generates n evenly spaced colors between a and b (inclusive).
func BlendN(a, b Color, n int) []Color {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []Color{a}
	}
	colors := make([]Color, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		colors[i] = Blend(a, b, t)
	}
	return colors
}

// Multiply performs multiply blend mode.
func Multiply(a, b Color) Color {
	return Color{
		R: uint8(uint16(a.R) * uint16(b.R) / 255),
		G: uint8(uint16(a.G) * uint16(b.G) / 255),
		B: uint8(uint16(a.B) * uint16(b.B) / 255),
	}
}

// Screen performs screen blend mode.
func Screen(a, b Color) Color {
	return Color{
		R: 255 - uint8(uint16(255-a.R)*uint16(255-b.R)/255),
		G: 255 - uint8(uint16(255-a.G)*uint16(255-b.G)/255),
		B: 255 - uint8(uint16(255-a.B)*uint16(255-b.B)/255),
	}
}

// Overlay combines multiply and screen based on the base color.
func Overlay(base, blend Color) Color {
	return Color{
		R: overlay8(base.R, blend.R),
		G: overlay8(base.G, blend.G),
		B: overlay8(base.B, blend.B),
	}
}

// Distance returns the Euclidean distance in RGB space.
func Distance(a, b Color) float64 {
	dr := float64(a.R) - float64(b.R)
	dg := float64(a.G) - float64(b.G)
	db := float64(a.B) - float64(b.B)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

func lerp8(a, b uint8, t float64) uint8 {
	return uint8(math.Round(float64(a)*(1-t) + float64(b)*t))
}

func overlay8(base, blend uint8) uint8 {
	if base < 128 {
		return uint8(2 * uint16(base) * uint16(blend) / 255)
	}
	return 255 - uint8(2*uint16(255-base)*uint16(255-blend)/255)
}
