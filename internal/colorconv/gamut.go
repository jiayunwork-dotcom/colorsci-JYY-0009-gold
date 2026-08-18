package colorconv

import "math"

// InLinearGamut reports whether the linear-light RGB triplet lies within
// the sRGB gamut (all channels within [0,1], small tolerance for float
// error).
func InLinearGamut(r, g, b float64) bool {
	const tol = 1e-7
	for _, v := range []float64{r, g, b} {
		if v < -tol || v > 1+tol {
			return false
		}
	}
	return true
}

// LabInGamut reports whether the CIELAB color can be represented as an
// sRGB8 color. The check converts to XYZ, then to linear RGB, then tests
// the linear triplet against the gamut bounds.
func LabInGamut(L, a, b float64) bool {
	x, y, z := LabToXYZ(L, a, b)
	r, g, bl := XYZToLinearRGB(x, y, z)
	return InLinearGamut(r, g, bl)
}

// LabInGamut8 reports whether the CIELAB color survives a round trip
// through 8-bit sRGB quantization without a change larger than tol in
// each Lab channel.
func LabInGamut8(L, a, b float64, tol float64) bool {
	c := LabToColor(L, a, b)
	L2, a2, b2 := ColorToLab(c)
	return math.Abs(L-L2) <= tol && math.Abs(a-a2) <= tol && math.Abs(b-b2) <= tol
}

// ClampLabToGamut returns the closest in-gamut CIELAB color by reducing
// chroma while preserving lightness and hue. Colors already in gamut are
// returned unchanged. The search is a bisection on the chroma scale
// (40 iterations, error below 2^-40 of the original chroma).
func ClampLabToGamut(L, a, b float64) (float64, float64, float64) {
	if LabInGamut(L, a, b) {
		return L, a, b
	}
	_, c, h := LabToLCh(L, a, b)
	lo, hi := 0.0, 1.0
	for i := 0; i < 40; i++ {
		mid := (lo + hi) / 2
		aL, aA, aB := LChToLab(L, c*mid, h)
		if LabInGamut(aL, aA, aB) {
			lo = mid
		} else {
			hi = mid
		}
	}
	return LChToLab(L, c*lo, h)
}
