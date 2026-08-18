// Package colorparse parses CSS color syntax into a normalized sRGB
// representation. It supports the functional notations rgb()/rgba() and
// hsl()/hsla() (both comma-separated legacy and space-separated modern
// forms), hex shorthand and longhand with optional alpha, and the W3C
// named color table.
//
// The zero value of Color is not meaningful; always use Parse or one of
// the typed constructors.
package colorparse

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Color is a normalized sRGB color with premultiplied-free semantics:
// R, G and B are in 8-bit precision and A is the alpha channel in the
// range [0,1]. Alpha is always stored explicitly; a fully opaque color
// has A == 1.
type Color struct {
	R uint8
	G uint8
	B uint8
	A float64
}

// Errors returned by Parse. Wrap with %w so errors.Is works.
var (
	ErrEmpty      = errors.New("colorparse: empty input")
	ErrBadSyntax  = errors.New("colorparse: bad color syntax")
	ErrBadChannel = errors.New("colorparse: channel out of range")
	ErrBadAlpha   = errors.New("colorparse: alpha out of range")
	ErrUnknown    = errors.New("colorparse: unknown color name")
	ErrHue        = errors.New("colorparse: bad hue")
)

// Parse interprets s as a CSS color and returns the normalized color.
// The input is case-insensitive. Recognized forms:
//
//	#rgb  #rgba  #rrggbb  #rrggbbaa
//	rgb(r g b)  rgb(r g b / a)  rgb(r,g,b)  rgb(r,g,b,a)
//	hsl(h s l)  hsl(h s l / a)  hsla(...)   (legacy comma forms)
//	<css-named-color>
//
// Any other input yields a wrapped ErrBadSyntax.
func Parse(s string) (Color, error) {
	t := strings.ToLower(strings.TrimSpace(s))
	if t == "" {
		return Color{}, fmt.Errorf("%w", ErrEmpty)
	}
	switch {
	case strings.HasPrefix(t, "#"):
		return ParseHex(t)
	case strings.HasPrefix(t, "rgb("), strings.HasPrefix(t, "rgba("):
		return ParseRGBFunc(t)
	case strings.HasPrefix(t, "hsl("), strings.HasPrefix(t, "hsla("):
		return ParseHSLFunc(t)
	}
	c, ok := named[t]
	if !ok {
		return Color{}, fmt.Errorf("%w: %q", ErrUnknown, s)
	}
	return c, nil
}

// ParseHex parses hex color syntax. Supported lengths (after '#'):
// 3 (rgb), 4 (rgba), 6 (rrggbb) and 8 (rrggbbaa). Shorthand digits are
// duplicated. The alpha byte maps to A in [0,1].
func ParseHex(s string) (Color, error) {
	hex := strings.TrimPrefix(s, "#")
	switch len(hex) {
	case 3, 4:
		if !allHex(hex) {
			return Color{}, fmt.Errorf("%w: bad hex digit in %q", ErrBadSyntax, s)
		}
		r, _ := hexByte(hex[0])
		g, _ := hexByte(hex[1])
		b, _ := hexByte(hex[2])
		a := byte(255)
		if len(hex) == 4 {
			a, _ = hexByte(hex[3])
		}
		return Color{R: r, G: g, B: b, A: float64(a) / 255}, nil
	case 6, 8:
		if !allHex(hex) {
			return Color{}, fmt.Errorf("%w: bad hex digit in %q", ErrBadSyntax, s)
		}
		r := pair(hex[0], hex[1])
		g := pair(hex[2], hex[3])
		b := pair(hex[4], hex[5])
		a := byte(255)
		if len(hex) == 8 {
			a = pair(hex[6], hex[7])
		}
		return Color{R: r, G: g, B: b, A: float64(a) / 255}, nil
	default:
		return Color{}, fmt.Errorf("%w: hex must be 3,4,6 or 8 digits, got %q", ErrBadSyntax, s)
	}
}

func allHex(s string) bool {
	for i := 0; i < len(s); i++ {
		if !isHex(s[i]) {
			return false
		}
	}
	return true
}

func isHex(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f')
}

func hexVal(b byte) uint8 {
	switch {
	case b >= '0' && b <= '9':
		return b - '0'
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10
	}
	return 0
}

func hexByte(b byte) (uint8, error) {
	return hexVal(b) * 17, nil
}

func pair(hi, lo byte) uint8 {
	return hexVal(hi)<<4 | hexVal(lo)
}

// ParseRGBFunc parses rgb()/rgba() functional syntax. Both the legacy
// comma-separated form and the modern space-separated form with an
// optional "/ alpha" are accepted. Channel values may be integers in
// [0,255] or percentages in [0%,100%]; mixing units in one call is an
// error. Alpha is a number in [0,1] or a percentage.
func ParseRGBFunc(s string) (Color, error) {
	body, ok := stripFunc(s, "rgb", "rgba")
	if !ok {
		return Color{}, fmt.Errorf("%w: not an rgb() function: %q", ErrBadSyntax, s)
	}
	args := splitArgs(body)

	channels := args
	alpha := 1.0
	if len(args) == 4 {
		channels = args[:3]
		a, err := alphaValue(args[3])
		if err != nil {
			return Color{}, err
		}
		alpha = a
	}

	r, unitR, err := channelValue(channels[0])
	if err != nil {
		return Color{}, err
	}
	g, unitG, err := channelValue(channels[1])
	if err != nil {
		return Color{}, err
	}
	b, unitB, err := channelValue(channels[2])
	if err != nil {
		return Color{}, err
	}
	if unitR != unitG || unitG != unitB {
		return Color{}, fmt.Errorf("%w: cannot mix number and percentage channels in %q", ErrBadSyntax, s)
	}
	return Color{R: clamp255(r), G: clamp255(g), B: clamp255(b), A: alpha}, nil
}

func splitArgs(body string) []string {
	body = strings.TrimSpace(body)
	if strings.Contains(body, "/") {
		parts := strings.SplitN(body, "/", 2)
		head := strings.Fields(parts[0])
		if len(head) != 3 {
			return []string{}
		}
		return []string{head[0], head[1], head[2], strings.TrimSpace(parts[1])}
	}
	return strings.Fields(strings.ReplaceAll(body, ",", " "))
}

// channelValue parses one RGB channel and returns its numeric value in
// 0..255 scale along with the unit kind used (percent or absolute).
func channelValue(tok string) (float64, int, error) {
	if strings.HasSuffix(tok, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(tok, "%"), 64)
		if err != nil {
			return 0, 0, fmt.Errorf("%w: bad channel %q", ErrBadSyntax, tok)
		}
		return v / 100 * 255, 1, nil
	}
	v, err := strconv.ParseFloat(tok, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: bad channel %q", ErrBadSyntax, tok)
	}
	return v, 0, nil
}

func alphaValue(tok string) (float64, error) {
	var v float64
	var err error
	if strings.HasSuffix(tok, "%") {
		v, err = strconv.ParseFloat(strings.TrimSuffix(tok, "%"), 64)
		if err != nil {
			return 0, fmt.Errorf("%w: bad alpha %q", ErrBadSyntax, tok)
		}
		v /= 100
	} else {
		v, err = strconv.ParseFloat(tok, 64)
		if err != nil {
			return 0, fmt.Errorf("%w: bad alpha %q", ErrBadSyntax, tok)
		}
	}
	if v < 0 || v > 1 {
		return 0, fmt.Errorf("%w: %v", ErrBadAlpha, v)
	}
	return v, nil
}

func clamp255(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(math.Round(v))
}

func stripFunc(s string, names ...string) (string, bool) {
	open := strings.IndexByte(s, '(')
	close := strings.LastIndexByte(s, ')')
	if open < 0 || close < open {
		return "", false
	}
	head := strings.TrimSpace(s[:open])
	for _, n := range names {
		if head == n {
			return s[open+1 : close], true
		}
	}
	return "", false
}

// ParseHSLFunc parses hsl()/hsla() functional syntax. The hue may carry
// an optional unit suffix: deg (default), rad, grad or turn. Saturation
// and lightness must be percentages. Legacy comma form and modern space
// form are both accepted.
func ParseHSLFunc(s string) (Color, error) {
	body, ok := stripFunc(s, "hsl", "hsla")
	if !ok {
		return Color{}, fmt.Errorf("%w: not an hsl() function: %q", ErrBadSyntax, s)
	}
	args := splitArgs(body)
	if len(args) < 3 {
		return Color{}, fmt.Errorf("%w: hsl() needs 3 channels: %q", ErrBadSyntax, s)
	}
	channels := args[:3]
	alpha := 1.0
	if len(args) == 4 {
		a, err := alphaValue(args[3])
		if err != nil {
			return Color{}, err
		}
		alpha = a
	}
	h, err := hueValue(channels[0])
	if err != nil {
		return Color{}, err
	}
	sat, err := percentChannel(channels[1], "saturation")
	if err != nil {
		return Color{}, err
	}
	light, err := percentChannel(channels[2], "lightness")
	if err != nil {
		return Color{}, err
	}
	r, g, b := hslToRGB(h, sat/100, light/100)
	return Color{R: clamp255(r * 255), G: clamp255(g * 255), B: clamp255(b * 255), A: alpha}, nil
}

func hueValue(tok string) (float64, error) {
	lower := strings.ToLower(tok)
	mult := 1.0
	switch {
	case strings.HasSuffix(lower, "deg"):
		lower = strings.TrimSuffix(lower, "deg")
	case strings.HasSuffix(lower, "rad"):
		lower = strings.TrimSuffix(lower, "rad")
		mult = 180 / math.Pi
	case strings.HasSuffix(lower, "grad"):
		lower = strings.TrimSuffix(lower, "grad")
		mult = 0.9
	case strings.HasSuffix(lower, "turn"):
		lower = strings.TrimSuffix(lower, "turn")
		mult = 360
	}
	v, err := strconv.ParseFloat(lower, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrHue, tok)
	}
	v *= mult
	v = math.Mod(v, 360)
	if v < 0 {
		v += 360
	}
	return v, nil
}

func percentChannel(tok, name string) (float64, error) {
	if !strings.HasSuffix(tok, "%") {
		return 0, fmt.Errorf("%w: %s must be a percentage, got %q", ErrBadSyntax, name, tok)
	}
	v, err := strconv.ParseFloat(strings.TrimSuffix(tok, "%"), 64)
	if err != nil {
		return 0, fmt.Errorf("%w: bad %s %q", ErrBadSyntax, name, tok)
	}
	if v < 0 || v > 100 {
		return 0, fmt.Errorf("%w: %s %v outside [0,100]", ErrBadChannel, name, v)
	}
	return v, nil
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

// hslToRGB converts hue (degrees, 0..360), saturation and lightness
// (both 0..1) to normalized sRGB channels in 0..1.
func hslToRGB(h, s, l float64) (float64, float64, float64) {
	if s == 0 {
		return l, l, l
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	hk := h / 360
	return hueToRGB(p, q, hk+1.0/3), hueToRGB(p, q, hk), hueToRGB(p, q, hk-1.0/3)
}

// Names returns the sorted list of supported named colors.
func Names() []string {
	out := make([]string, 0, len(named))
	for k := range named {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
