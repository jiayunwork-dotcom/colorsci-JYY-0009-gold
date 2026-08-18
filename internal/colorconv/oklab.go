package colorconv

import (
	"math"

	"colorsci/internal/colorparse"
)

// OKLab is a perceptually uniform color space designed by Björn Ottosson.
// It improves upon CIELAB with better hue linearity.

// RGBToOKLab converts sRGB to OKLab.
func RGBToOKLab(c colorparse.Color) (float64, float64, float64) {
	r := Linearize(float64(c.R) / 255)
	g := Linearize(float64(c.G) / 255)
	b := Linearize(float64(c.B) / 255)

	l := 0.4122214708*r + 0.5363325363*g + 0.0514459929*b
	m := 0.2119034982*r + 0.6806995451*g + 0.1073969566*b
	s := 0.0883024619*r + 0.2817188376*g + 0.6299787005*b

	l_ := math.Cbrt(l)
	m_ := math.Cbrt(m)
	s_ := math.Cbrt(s)

	L := 0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_
	A := 1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_
	B := 0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_

	return L, A, B
}

// OKLabToRGB converts OKLab to sRGB.
func OKLabToRGB(L, A, B float64) colorparse.Color {
	l_ := L + 0.3963377774*A + 0.2158037573*B
	m_ := L - 0.1055613458*A - 0.0638541728*B
	s_ := L - 0.0894841775*A - 1.2914855480*B

	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_

	r := +4.0767416621*l - 3.3077115913*m + 0.2309699292*s
	g := -1.2684380046*l + 2.6097574011*m - 0.3413193965*s
	b := -0.0041960863*l - 0.7034186147*m + 1.7076147010*s

	return colorparse.Color{
		R: clamp8(Encode(r)),
		G: clamp8(Encode(g)),
		B: clamp8(Encode(b)),
	}
}

// OKLabToLCh converts OKLab to cylindrical LCh form.
func OKLabToLCh(L, A, B float64) (float64, float64, float64) {
	C := math.Sqrt(A*A + B*B)
	h := math.Atan2(B, A) * 180 / math.Pi
	if h < 0 {
		h += 360
	}
	return L, C, h
}

// OKLChToLab converts OKLCh back to OKLab.
func OKLChToLab(L, C, h float64) (float64, float64, float64) {
	hRad := h * math.Pi / 180
	A := C * math.Cos(hRad)
	B := C * math.Sin(hRad)
	return L, A, B
}

func clamp8(v float64) uint8 {
	x := v * 255
	if x < 0 {
		return 0
	}
	if x > 255 {
		return 255
	}
	return uint8(math.Round(x))
}
