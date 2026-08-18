// Package colormetric computes color difference metrics (CIE76, CIE94,
// CIEDE2000) and WCAG contrast ratios for accessibility checks. All
// metrics operate on CIELAB or sRGB8 colors supplied by the caller.
package colormetric

import (
	"math"

	"colorsci/internal/colorconv"
	"colorsci/internal/colorparse"
)

// Lab is a CIELAB color triplet.
type Lab struct {
	L float64
	A float64
	B float64
}

// LabOf returns the Lab triplet for an sRGB8 color.
func LabOf(c colorparse.Color) Lab {
	l, a, b := colorconv.ColorToLab(c)
	return Lab{L: l, A: a, B: b}
}

// DeltaE76 returns the CIE76 color difference (Euclidean distance in
// Lab). It is symmetric.
func DeltaE76(x, y Lab) float64 {
	dl := x.L - y.L
	da := x.A - y.A
	db := x.B - y.B
	return math.Sqrt(dl*dl + da*da + db*db)
}

// DeltaE94 returns the CIE94 color difference using the graphic-arts
// weights (kL=1, K1=0.045, K2=0.015). The metric is not symmetric: the
// reference color is x, the sample is y.
func DeltaE94(x, y Lab) float64 {
	c1 := math.Hypot(x.A, x.B)
	c2 := math.Hypot(y.A, y.B)
	dl := x.L - y.L
	dc := c1 - c2
	dh2 := (x.A-y.A)*(x.A-y.A) + (x.B-y.B)*(x.B-y.B) - dc*dc
	if dh2 < 0 {
		dh2 = 0
	}
	sl := 1.0
	sc := 1 + 0.045*c1
	sh := 1 + 0.015*c1
	return math.Sqrt((dl/sl)*(dl/sl) + (dc/sc)*(dc/sc) + dh2/(sh*sh))
}

// DeltaE2000 returns the CIEDE2000 color difference. It is symmetric
// and uses the standard parametric weights kL=kC=kH=1. The result is
// zero for identical colors.
func DeltaE2000(x, y Lab) float64 {
	// Chroma of each color in Lab (a' adjustments).
	c1 := math.Hypot(x.A, x.B)
	c2 := math.Hypot(y.A, y.B)

	// Mean chroma and the G factor (0.5 * (1 - sqrt(Cbar^7/(Cbar^7+25^7)))).
	cbar := (c1 + c2) / 2
	g := 0.5 * (1 - math.Sqrt(math.Pow(cbar, 7)/(math.Pow(cbar, 7)+math.Pow(25, 7))))

	a1 := (1 + g) * x.A
	a2 := (1 + g) * y.A
	cp1 := math.Hypot(a1, x.B)
	cp2 := math.Hypot(a2, y.B)

	h1 := hueDeg(a1, x.B)
	h2 := hueDeg(a2, y.B)

	dl := y.L - x.L
	dc := cp2 - cp1
	dh := deltaHue(h1, h2, cp1, cp2)
	dh2 := 2 * math.Sqrt(cp1*cp2) * math.Sin(dh/2*math.Pi/180)

	lbar := (x.L + y.L) / 2
	cpbar := (cp1 + cp2) / 2
	hbar := meanHue(h1, h2, cp1, cp2)

	t := 1 - 0.17*math.Cos((hbar-30)*math.Pi/180) +
		0.24*math.Cos((2*hbar)*math.Pi/180) +
		0.32*math.Cos((3*hbar+6)*math.Pi/180) -
		0.20*math.Cos((4*hbar-63)*math.Pi/180)

	sl := 1 + 0.015*math.Pow(lbar-50, 2)/math.Sqrt(20+math.Pow(lbar-50, 2))
	sc := 1 + 0.045*cpbar
	sh := 1 + 0.015*cpbar*t

	rt := -2 * math.Sqrt(math.Pow(cpbar, 7)/(math.Pow(cpbar, 7)+math.Pow(25, 7))) *
		math.Sin(2*(30*math.Exp(-math.Pow((hbar-275)/25, 2)))*math.Pi/180)

	term1 := dl / sl
	term2 := dc / sc
	term3 := dh2 / sh
	return math.Sqrt(term1*term1 + term2*term2 + term3*term3 + rt*term2*term3)
}

// hueDeg returns the hue angle in degrees in [0,360) for a Lab color.
// Achromatic colors (radius < 1e-9) report 0.
func hueDeg(a, b float64) float64 {
	r := math.Hypot(a, b)
	if r < 1e-9 {
		return 0
	}
	h := math.Atan2(b, a) * 180 / math.Pi
	if h < 0 {
		h += 360
	}
	return h
}

func deltaHue(h1, h2, c1, c2 float64) float64 {
	if c1*c2 == 0 {
		return 0
	}
	d := h2 - h1
	if math.Abs(d) <= 180 {
		return d
	}
	if d > 180 {
		return d - 360
	}
	return d + 360
}

func meanHue(h1, h2, c1, c2 float64) float64 {
	if c1*c2 == 0 {
		return h1 + h2
	}
	d := math.Abs(h1 - h2)
	if d <= 180 {
		return (h1 + h2) / 2
	}
	if h1+h2 < 360 {
		return (h1 + h2 + 360) / 2
	}
	return (h1 + h2 - 360) / 2
}
