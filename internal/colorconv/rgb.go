// Package colorconv converts colors between the sRGB, linear RGB, CIE
// XYZ (D65), CIELAB and CIELCh color spaces, and provides gamut
// clamping. All conversions operate on normalized values: sRGB channels
// in 0..255 (from colorparse.Color), linear/XYZ/Lab in their native
// ranges.
package colorconv

import (
	"math"

	"colorsci/internal/colorparse"
)

// D65 reference white in XYZ for sRGB.
const (
	WhiteX = 0.95047
	WhiteY = 1.0
	WhiteZ = 1.08883
)

// Linearize applies the sRGB transfer function to a channel value in
// [0,1], returning the linear-light intensity.
func Linearize(v float64) float64 {
	v = math.Max(0, math.Min(1, v))
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// Encode applies the inverse sRGB transfer function to a linear-light
// intensity, returning the gamma-encoded channel in [0,1].
func Encode(v float64) float64 {
	if v <= 0.0031308 {
		return 12.92 * v
	}
	return 1.055*math.Pow(math.Max(v, 0), 1/2.4) - 0.055
}

// RGBToLinear converts 8-bit sRGB channels to linear-light values in
// [0,1].
func RGBToLinear(r, g, b uint8) (float64, float64, float64) {
	return Linearize(float64(r) / 255),
		Linearize(float64(g) / 255),
		Linearize(float64(b) / 255)
}

// LinearToRGB8 converts linear-light channels in [0,1] back to 8-bit
// sRGB, clamping out-of-range values and rounding.
func LinearToRGB8(r, g, b float64) (uint8, uint8, uint8) {
	enc := func(v float64) uint8 {
		v = math.Max(0, math.Min(1, v))
		return uint8(math.Round(Encode(v) * 255))
	}
	return enc(r), enc(g), enc(b)
}

// LinearToRGB8From returns the 8-bit sRGB color for linear channels.
func LinearToRGB8From(r, g, b float64) colorparse.Color {
	r8, g8, b8 := LinearToRGB8(r, g, b)
	return colorparse.Color{R: r8, G: g8, B: b8, A: 1}
}

// sRGB D65 matrices (IEC 61966-2-1).
const (
	m00 = 0.4124564
	m01 = 0.3575761
	m02 = 0.1804375
	m10 = 0.2126729
	m11 = 0.7151522
	m12 = 0.0721750
	m20 = 0.0193339
	m21 = 0.1191920
	m22 = 0.9503041

	inv00 = 3.2404542
	inv01 = -1.5371385
	inv02 = -0.4985314
	inv10 = -0.9692660
	inv11 = 1.8760108
	inv12 = 0.0415560
	inv20 = 0.0556434
	inv21 = -0.2040259
	inv22 = 1.0572252
)

// LinearRGBToXYZ converts linear-light RGB (D65) to CIE XYZ.
func LinearRGBToXYZ(r, g, b float64) (float64, float64, float64) {
	x := m00*r + m01*g + m02*b
	y := m10*r + m11*g + m12*b
	z := m20*r + m21*g + m22*b
	return x, y, z
}

// XYZToLinearRGB converts CIE XYZ to linear-light RGB (D65). Results
// may fall outside [0,1] for out-of-gamut colors.
func XYZToLinearRGB(x, y, z float64) (float64, float64, float64) {
	r := inv00*x + inv01*y + inv02*z
	g := inv10*x + inv11*y + inv12*z
	b := inv20*x + inv21*y + inv22*z
	return r, g, b
}

// ColorToXYZ converts a colorparse.Color (sRGB8 + alpha, alpha is
// ignored) to CIE XYZ (D65).
func ColorToXYZ(c colorparse.Color) (float64, float64, float64) {
	r, g, b := RGBToLinear(c.R, c.G, c.B)
	return LinearRGBToXYZ(r, g, b)
}

// XYZToColor converts CIE XYZ to an 8-bit sRGB color, clamping
// out-of-gamut values.
func XYZToColor(x, y, z float64) colorparse.Color {
	r, g, b := XYZToLinearRGB(x, y, z)
	return LinearToRGB8From(r, g, b)
}
