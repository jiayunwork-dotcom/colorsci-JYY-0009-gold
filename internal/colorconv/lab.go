package colorconv

import (
	"math"

	"colorsci/internal/colorparse"
)

// CIELAB conversion constants (CIE 1976).
const (
	eps   = 216.0 / 24389.0
	kappa = 24389.0 / 27.0
)

func labF(t float64) float64 {
	if t > eps {
		return math.Cbrt(t)
	}
	return (kappa*t + 16) / 116
}

func labInvF(t float64) float64 {
	t3 := t * t * t
	if t3 > eps {
		return t3
	}
	return (116*t - 16) / kappa
}

// XYZToLab converts CIE XYZ (D65) to CIELAB. L is in [0,100], a and b
// are unbounded in theory.
func XYZToLab(x, y, z float64) (float64, float64, float64) {
	fx := labF(x / WhiteX)
	fy := labF(y / WhiteY)
	fz := labF(z / WhiteZ)
	L := 116*fy - 16
	a := 500 * (fx - fy)
	b := 200 * (fy - fz)
	return L, a, b
}

// LabToXYZ converts CIELAB back to CIE XYZ (D65).
func LabToXYZ(L, a, b float64) (float64, float64, float64) {
	fy := (L + 16) / 116
	fx := fy + a/500
	fz := fy - b/200
	x := labInvF(fx) * WhiteX
	y := labInvF(fy) * WhiteY
	z := labInvF(fz) * WhiteZ
	return x, y, z
}

// LabToLCh converts CIELAB to CIELCh. Hue h is in degrees in [0,360);
// chroma C is the radial distance from the L axis. Achromatic colors
// (C < 1e-9) report h = 0.
func LabToLCh(L, a, b float64) (float64, float64, float64) {
	c := math.Hypot(a, b)
	h := 0.0
	if c > 1e-9 {
		h = math.Atan2(b, a) * 180 / math.Pi
		if h < 0 {
			h += 360
		}
	}
	return L, c, h
}

// LChToLab converts CIELCh back to CIELAB.
func LChToLab(L, c, h float64) (float64, float64, float64) {
	hr := h * math.Pi / 180
	return L, c * math.Cos(hr), c * math.Sin(hr)
}

// ColorToLab converts an sRGB8 color to CIELAB (D65), ignoring alpha.
func ColorToLab(c colorparse.Color) (float64, float64, float64) {
	x, y, z := ColorToXYZ(c)
	return XYZToLab(x, y, z)
}

// LabToColor converts CIELAB to an 8-bit sRGB color, clamping
// out-of-gamut values.
func LabToColor(L, a, b float64) colorparse.Color {
	x, y, z := LabToXYZ(L, a, b)
	return XYZToColor(x, y, z)
}

// RGBToHSL converts 8-bit sRGB channels to HSL. H is in [0,360), S and
// L are in [0,1].
func RGBToHSL(r, g, b uint8) (float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255
	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	l := (max + min) / 2
	if max == min {
		return 0, 0, l
	}
	d := max - min
	s := d / (1 - math.Abs(2*l-1))
	var h float64
	switch max {
	case rf:
		h = 60 * math.Mod((gf-bf)/d, 6)
	case gf:
		h = 60 * ((bf-rf)/d + 2)
	default:
		h = 60 * ((rf-gf)/d + 4)
	}
	if h < 0 {
		h += 360
	}
	return h, s, l
}

// HSLToRGB8 converts hue (degrees), saturation and lightness (both
// [0,1]) to 8-bit sRGB channels.
func HSLToRGB8(h, s, l float64) (uint8, uint8, uint8) {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	hk := h / 360
	if s == 0 {
		v := uint8(math.Round(math.Max(0, math.Min(1, l)) * 255))
		return v, v, v
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	r := hueToRGB(p, q, hk+1.0/3)
	g := hueToRGB(p, q, hk)
	b := hueToRGB(p, q, hk-1.0/3)
	clamp := func(v float64) uint8 {
		v = math.Max(0, math.Min(1, v))
		return uint8(math.Round(v * 255))
	}
	return clamp(r), clamp(g), clamp(b)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t < 1.0/2:
		return q
	case t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}
