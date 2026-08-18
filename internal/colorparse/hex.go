package colorparse

import (
	"fmt"
	"math/rand"
	"strings"
)

// NormalizeHex takes a hex color string and returns a normalized 6-digit
// lowercase hex with '#' prefix. Shorthand (#rgb) is expanded. Returns
// an error if the input is not valid hex syntax.
func NormalizeHex(s string) (string, error) {
	c, err := ParseHex(s)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B), nil
}

// RandomHex generates a random fully-opaque hex color string (#rrggbb)
// using the provided random source. If rng is nil, a default source is used.
func RandomHex(rng *rand.Rand) string {
	var r, g, b uint8
	if rng != nil {
		r = uint8(rng.Intn(256))
		g = uint8(rng.Intn(256))
		b = uint8(rng.Intn(256))
	} else {
		r = uint8(rand.Intn(256))
		g = uint8(rand.Intn(256))
		b = uint8(rand.Intn(256))
	}
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// IsValidHex reports whether s is a syntactically valid CSS hex color.
func IsValidHex(s string) bool {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(s) != 3 && len(s) != 4 && len(s) != 6 && len(s) != 8 {
		return false
	}
	return allHex(strings.ToLower(s))
}

// HexToUpper returns the hex color string with uppercase hex digits.
func HexToUpper(s string) string {
	if !strings.HasPrefix(s, "#") {
		return "#" + strings.ToUpper(s)
	}
	return "#" + strings.ToUpper(s[1:])
}
